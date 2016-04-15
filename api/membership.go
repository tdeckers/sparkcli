package api

import (
	"errors"
	"github.com/tdeckers/sparkcli/util"
	"net/url"
)

type MemberService struct {
	Client *util.Client
}

type Membership struct {
	Id                string `json:"id,omitempty"`
	RoomId            string `json:"roomId,omitempty"`
	PersonId          string `json:"personId,omitempty"`
	PersonEmail       string `json:"personEmail,omitempty"`
	PersonDisplayName string `json:"personDisplayName,omitempty"`
	IsModerator       bool   `json:"isModerator,omitempty"`
	IsMonitor         bool   `json:"isMonitor,omitempty"`
	Created           string `json:"created,omitempty"`
}

type MembershipItems struct {
	Items []Membership `json:"items"`
}

func (m MemberService) List(roomId string, personId string, personEmail string) (*[]Membership, error) {
	v := url.Values{}
	if roomId != "" {
		v.Add("roomId", roomId)
	}
	if personId != "" {
		v.Add("personId", personId)
	}
	if personEmail != "" {
		v.Add("personEmail", personEmail)
	}
	req, err := m.Client.NewGetRequest("/memberships?" + v.Encode())
	if err != nil {
		return nil, err
	}
	var result MembershipItems
	_, err = m.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result.Items, nil
}

func (m MemberService) Create(roomId, personId, personEmail string) (*Membership, error) {
	// check default room id
	config := util.GetConfiguration()
	if roomId == "-" {
		if config.DefaultRoomId != "" {
			roomId = config.DefaultRoomId
		} else {
			return nil, errors.New("No DefaultRoomId configured.")
		}
	}
	ms := Membership{RoomId: roomId, PersonId: personId, PersonEmail: personEmail}
	req, err := m.Client.NewPostRequest("/memberships", ms)
	if err != nil {
		return nil, err
	}
	var result Membership
	_, err = m.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m MemberService) Get(id string) (*Membership, error) {
	req, err := m.Client.NewGetRequest("/memberships/" + id)
	if err != nil {
		return nil, err
	}
	var result Membership
	_, err = m.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m MemberService) Update(id string, isModerator bool) (*Membership, error) {
	ms := Membership{IsModerator: isModerator}
	req, err := m.Client.NewPutRequest("/memberships/"+id, ms)
	if err != nil {
		return nil, err
	}
	var result Membership
	_, err = m.Client.Do(req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m MemberService) Delete(id string) error {
	req, err := m.Client.NewDeleteRequest("/memberships/" + id)
	if err != nil {
		return err
	}
	_, err = m.Client.Do(req, nil)
	if err != nil {
		return err
	}
	return nil
}
