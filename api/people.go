package api

import (
	"github.com/tdeckers/sparkcli/util"
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

func (p PeopleService) GetMe() (*People, error) {
	req, err := p.Client.NewGetRequest("/people/me")
	if err != nil {
		return nil, err
	}
	var result People
	res, err := p.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	err = util.CheckStatusOk(res)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
