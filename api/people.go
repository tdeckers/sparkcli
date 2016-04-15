package api

import (
	"errors"
	"github.com/tdeckers/sparkcli/util"
	"net/url"
)

type PeopleService struct {
	Client *util.Client
}

type People struct {
	Id          string   `json:"id,omitempty"`
	Emails      []string `json:"emails,omitempty"`
	DisplayName string   `json:"displayName,omitempty"`
	Avatar      string   `json:"avatar,omitempty"`
	Created     string   `json:"created,omitempty"`
}

type PeopleItems struct {
	Items []People `json:"items"`
}

func (p PeopleService) List(email string, displayName string) (*[]People, error) {
	if email == "" && displayName == "" {
		// TODO: don't need to create this message.  Just return what service returns.
		//{
		//	"message": "Email or displayName should be specified.",
		//	"errors": [
		//		{
		//			"description": "Email or displayName should be specified."
		//		}
		//	],
		//	"trackingId": "NA_4de291c7-f857-4c3b-a02d-5129e7cea02c"
		//}
		return nil, errors.New("Email or displayName should be specified")
	}
	v := url.Values{}
	if email != "" {
		v.Add("email", email)
	}
	// TODO: check searching by name, doesn't seem to always work?!
	if displayName != "" {
		v.Add("displayName", displayName)
	}
	req, err := p.Client.NewGetRequest("/people?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result PeopleItems
	_, err = p.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Items, nil
}

func (p PeopleService) Get(id string) (*People, error) {
	req, err := p.Client.NewGetRequest("/people/" + id)
	if err != nil {
		return nil, err
	}
	var result People
	_, err = p.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (p PeopleService) GetMe() (*People, error) {
	return p.Get("me")
}
