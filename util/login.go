package util

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Login struct {
	// TODO: should this be a pointer (reference)?
	config *Configuration
	client *Client
}

type Tokens struct {
	AccessToken    string  `json:"access_token"`
	AccessExpires  float64 `json:"expires_in"`
	RefreshToken   string  `json:"refresh_token"`
	RefreshExpires float64 `json:"refresh_token_expires_in"`
}

func NewLogin(config *Configuration, client *Client) Login {
	return Login{config: config, client: client}
}

func (l Login) Authorize() {
	// Check if credentials are set.
	l.config.checkConfAuth()

	// if tokens exist, try them.
	if l.config.AccessToken != "" && l.config.RefreshToken != "" {
		err := l.test()
		// TODO: what is another error than 401 pops up? Find way to detect 401 here.
		if err == nil {
			log.Println("Already logged in.")
			return
		}
		log.Printf("Error: %s", err)

	}

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

	log.Printf("Access token: %s", tokens.AccessToken)
	log.Printf("Refresh token: %s", tokens.RefreshToken)

	l.storeToken(tokens)
}

func (l Login) RefreshToken() {
	log.Print("Refreshing token...")
	// Post form to obtain access token based on refresh token (OAuth)
	res, err := http.PostForm(l.config.BaseUrl+"/access_token",
		url.Values{"grant_type": {"refresh_token"},
			"client_id":     {l.config.ClientId},
			"client_secret": {l.config.ClientSecret},
			"refresh_token": {l.config.AuthCode}})
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

	l.storeToken(tokens)
}

func (l Login) storeToken(tokens *Tokens) {

	// http://blog.golang.org/json-and-go#TOC_5.
	l.config.AccessToken = tokens.AccessToken
	// typically 14 days
	l.config.AccessExpires = tokens.AccessExpires
	l.config.RefreshToken = tokens.RefreshToken
	// typically 90 days
	l.config.RefreshExpires = tokens.RefreshExpires
	log.Println("Saving config")
	l.config.Save()

}

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
