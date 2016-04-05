package api

import (
	"github.com/tdeckers/sparkcli/util"
)

type RoomService struct {
	Client *util.Client
}

type Room struct {
	Id           string `json:"id,omitempty"`
	Title        string `json:"title"`
	Created      string `json:"created,omitempty"`
	LastActivity string `json:"lastActivity,omitempty"`
	IsLocked     bool   `json:"isLocked,omitempty"`
}

type RoomItems struct {
	Items []Room `json:"items"`
}

func (r RoomService) List() (*[]Room, error) {
	req, err := r.Client.NewGetRequest("/rooms")
	if err != nil {
		return nil, err
	}
	var result RoomItems
	res, err := r.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	err = util.CheckStatusOk(res)
	if err != nil {
		return nil, err
	}
	return &result.Items, nil
}

func (r RoomService) Create(name string) (*Room, error) {
	room := Room{Title: name}
	req, err := r.Client.NewPostRequest("/rooms", room)
	if err != nil {
		return nil, err
	}
	var result Room
	res, err := r.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	err = util.CheckStatusOk(res)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r RoomService) Get(id string) (*Room, error) {
	req, err := r.Client.NewGetRequest("/rooms/" + id)
	if err != nil {
		return nil, err
	}
	var result Room
	res, err := r.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	err = util.CheckStatusOk(res)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// TODO: add Update

func (r RoomService) Delete(id string) error {
	req, err := r.Client.NewDeleteRequest("/rooms/" + id)
	if err != nil {
		return err
	}
	res, err := r.Client.Do(req, nil)
	if err != nil {
		return err
	}
	err = util.CheckStatusOk(res)
	return err
}
