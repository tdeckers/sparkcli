package util

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
)

// Login allows authorization against the Cisco Spark service.
// It supports both bot and integration accounts.
type Login struct {
	config *Configuration
	client *Client
}

// Tokens is used to parse responses from access token requests
type Tokens struct {
	AccessToken    string  `json:"access_token"`
	AccessExpires  float64 `json:"expires_in"`
	RefreshToken   string  `json:"refresh_token"`
	RefreshExpires float64 `json:"refresh_token_expires_in"`
}

// NewLogin creates a new Login and initializes it.
func NewLogin(config *Configuration, client *Client) Login {
	return Login{config: config, client: client}
}

// Authorize will verify is proper a proper access token is available.  If not
// it will attempt to use the OAuth integration flow. to obtain an access token
// based on the provided ClientId, ClientSecret and AuthCode in the
// configuration.
func (l Login) Authorize() {
	// Check if AccessToken is present
	tokenPresent := l.config.checkAccessToken()
	if tokenPresent {
		// Verify if token works.
		err := l.test()
		if err != nil {
			l.loginAsIntegration()
		} else { // Success!
			return
		}
	} else { // AccessToken not present
		l.loginAsIntegration()
	}
}

// loginAsIntegration implements the OAuth grant flow for integration accouns.
// it expects a configuration file to be available with ClientId, ClientSecret
// and AuthCode set.
// On successful authentication it will store the AccessToken and RefreshToken
// in the configuration file for further use.  On failure it will exit the
// program.
func (l Login) loginAsIntegration() {
	// Check if client credentials are set.
	err := l.config.checkClientConfig()
	if err != nil { // If client credentials are not set...
		log.Fatalf("Not configured properly: %s", err)
	}
	// client credentials properly set, let's continue.

	log.Println("Authorizing...")
	// Post form to obtain access token based on authorization code (OAuth)
	res, err := http.PostForm(l.config.BaseUrl+"/access_token",
		url.Values{"grant_type": {"authorization_code"},
			"client_id":     {l.config.ClientId},
			"client_secret": {l.config.ClientSecret},
			"code":          {l.config.AuthCode},
			"redirect_uri":  {l.config.RedirectUri}})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// if 401, reauthorize? or refresh key.
	if res.StatusCode == 401 {
		log.Print("Unauthorized (401) - trying to refresh token")
		l.RefreshToken()
	} else if res.StatusCode != 200 {
		log.Fatal("Unexpected status code ", res.StatusCode)
	}

	// Parse json code into Tokens struct
	decoder := json.NewDecoder(res.Body)
	tokens := new(Tokens)
	err = decoder.Decode(&tokens)
	if err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	log.Printf("Access token: %s", tokens.AccessToken)
	log.Printf("Refresh token: %s", tokens.RefreshToken)

	l.storeToken(tokens, false)
}

// RefreshToken uses the ClientId, ClientSecret and RefreshToken from the
// configuration file and attempt to obtain a new access token.
// On success, the new AccessToken is written into the configuration
// file.  The RefreshToken remains the same, its expiry is reset.
// Note that sparkcli doesn't track token expiry.
func (l Login) RefreshToken() {
	log.Print("Refreshing token...")
	// Post form to obtain access token based on refresh token (OAuth)
	res, err := http.PostForm(l.config.BaseUrl+"/access_token",
		url.Values{"grant_type": {"refresh_token"},
			"client_id":     {l.config.ClientId},
			"client_secret": {l.config.ClientSecret},
			"refresh_token": {l.config.RefreshToken}})
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	// if 401, reauthorize?
	if res.StatusCode == 401 {
		log.Print("Unauthorized (401)")
		l.config.PrintAuthUrl()
		os.Exit(1)
	} else if res.StatusCode != 200 {
		log.Fatal("Unexpected status code ", res.StatusCode)
	}

	// Parse json code into Tokens struct
	decoder := json.NewDecoder(res.Body)
	tokens := new(Tokens)
	err = decoder.Decode(&tokens)
	if err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	l.storeToken(tokens, true)

	log.Printf("Successfully refreshed token.")
}

// storeToken writes tokens to the configuration file.  When refresh
// is true, it will not overwrite RefreshToken and RefreshExpires (since
// these will be empty during refresh)
func (l Login) storeToken(tokens *Tokens, refresh bool) {

	// http://blog.golang.org/json-and-go#TOC_5.
	l.config.AccessToken = tokens.AccessToken
	// typically 14 days
	l.config.AccessExpires = tokens.AccessExpires
	// A refresh doesn't repeat the refresh token, so let's not
	// overwrite with an empty value here!
	if !refresh {
		l.config.RefreshToken = tokens.RefreshToken
		// typically 90 days
		l.config.RefreshExpires = tokens.RefreshExpires
	}
	log.Println("Saving config")
	l.config.Save()

}

// test access to the Cisco Spark service to ensure authentication works as
// expected.  Returns an error if the service request fails.
func (l Login) test() error {
	req, err := l.client.NewGetRequest("/people/me")
	if err != nil {
		log.Fatalf("Error testing connection: %s", err)
	}
	var result interface{}
	res, err := l.client.Do(req, &result)
	if err != nil {
		log.Fatalf("Error testing connection: %s", err)
	}
	if res.StatusCode == 401 {
		return errors.New("401 Unauthorized")
	}
	if res.StatusCode != 200 {
		// TODO: what should we do in case of another error while testing?
		log.Printf("Got response code %v while testing.", res.StatusCode)
	}
	return nil
}
