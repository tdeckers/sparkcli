package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

const (
	userAgent = "spark-cli"
)

type Client struct {
	client *http.Client

	userAgent string

	config *Configuration
}

func NewClient(config *Configuration) *Client {
	c := &Client{client: http.DefaultClient, userAgent: userAgent, config: config}
	return c
}

func (c *Client) NewRequest(method string, path string, body interface{}) (*http.Request, error) {
	// concat base url and request url
	reqUrl, err := url.Parse(c.config.BaseUrl + path)
	if err != nil {
		return nil, err
	}
	var bodyBuffer *bytes.Buffer
	var req *http.Request
	// if body is present (likely for POST), then marshal and create buffer
	if body != nil {
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyBuffer = bytes.NewBuffer(bodyJson)
		log.Printf("Sending: %s", bodyBuffer)
		// Create request with body
		req, err = http.NewRequest(method, reqUrl.String(), bodyBuffer)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		// Create request without body
		req, err = http.NewRequest(method, reqUrl.String(), nil)
		if err != nil {
			return nil, err
		}
	}
	// Add other headers (that apply to all requests)
	req.Header.Set("Authorization", "Bearer "+c.config.AccessToken)
	return req, nil
}

func (c *Client) NewGetRequest(path string) (*http.Request, error) {
	return c.NewRequest("GET", path, nil)
}

func (c *Client) NewPostRequest(path string, body interface{}) (*http.Request, error) {
	return c.NewRequest("POST", path, body)
}

func (c *Client) NewDeleteRequest(path string) (*http.Request, error) {
	return c.NewRequest("DELETE", path, nil)
}

func (c *Client) Do(req *http.Request, to interface{}) (*http.Response, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// If 401, let's try to refresh tokens and try again.
	if res.StatusCode == 401 {
		login := Login{config: c.config, client: c}
		login.RefreshToken()
		res, err = c.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
	}
	//	log.Printf("Response: %s", res.Status)
	// TODO: check response code.  If 401, handle it!  See spark.go
	// other errors should be handed downstream (incl. second 401)
	if to != nil {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(&to)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// error if status code is not in 2XX range
func CheckStatusOk(res *http.Response) error {
	if 200 < res.StatusCode && res.StatusCode > 299 {
		// TODO: if 400 message, parse body and get actual message
		// {
		//	"message": "Failed to create room.",
		//	"errors": [
		//		{
		//			"description": "Failed to create room."
		//		}
		//	],
		//	"trackingId": "NA_f6e19aac-3a72-46d2-88ec-643f4d12fcbd"
		//}
		return errors.New(res.Status)
	}
	return nil
}
