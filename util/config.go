package util

import (
	"bufio"
	"bytes"
	"github.com/BurntSushi/toml"
	"log"
	"net/url"
	"os"
)

const (
	redirectUrl = "http://files.ducbase.com/code.html"
	scope       = "spark:people_read spark:rooms_read spark:rooms_write " +
		"spark:messages_read spark:messages_write"
	baseUrl = "https://api.ciscospark.com/v1"
)

type Configuration struct {
	BaseUrl        string
	ClientId       string
	ClientSecret   string
	AuthCode       string
	RedirectUri    string
	Scope          string
	AccessToken    string
	AccessExpires  float64
	RefreshToken   string
	RefreshExpires float64
}

func (c *Configuration) Load() {
	//TODO: check if empty after loading, else initalize
	if c.RedirectUri == "" {
		c.RedirectUri = redirectUrl
	}
	if c.Scope == "" {
		c.Scope = scope
	}
	if c.BaseUrl == "" {
		c.BaseUrl = baseUrl
	}

	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Fatalln("Failed to open file", err)
		return
	}
	log.Println("File loaded for " + c.ClientId)
}

func (c Configuration) save() {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		log.Fatalln("Failed to encode config", err)
	}
	f, err := os.Create("config.toml")
	if err != nil {
		log.Fatalln("Failed to create file", err)
		return
	}

	w := bufio.NewWriter(f)
	buf.WriteTo(w)
	w.Flush()
}

func (c Configuration) checkConfAuth() {
	if c.ClientId == "" {
		log.Fatalln("ClientId not configured")
	}
	if c.ClientSecret == "" {
		log.Fatalln("ClientSecret not configured")
	}
	if c.AuthCode == "" {
		c.PrintAuthUrl()
		log.Fatalln("AuthCode not configured")
	}
}

func (c Configuration) PrintAuthUrl() {
	log.Printf("Visit \n%s/authorize?%s",
		c.BaseUrl,
		url.Values{"response_type": {"code"},
			"client_id":    {c.ClientId},
			"redirect_uri": {c.RedirectUri},
			"scope":        {c.Scope}}.Encode(),
	)

}
