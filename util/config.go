package util

import (
	"bufio"
	"bytes"
	"github.com/BurntSushi/toml"
	"log"
	"net/url"
	"os"
	"os/user"
	"errors"
)

const (
	redirectUrl = "http://files.ducbase.com/code.html"
	scope       = "spark:people_read spark:rooms_read spark:rooms_write " +
		"spark:messages_read spark:messages_write spark:memberships_read " +
		"spark:memberships_write"
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
	DefaultRoomId  string
}

var configFile string

func init() {
	configFile = findConfigFile()
}

// golang singletons: http://marcio.io/2015/07/singleton-pattern-in-go/
var instance *Configuration

func GetConfiguration() *Configuration {
	if instance == nil {
		instance = &Configuration{}
	}
	return instance
}

func (c *Configuration) Load() {
	// TODO:change this to log to stderr, actuall all logs to stderr?
	//log.Printf("Using configuration at %s\n", configFile)

	if _, err := toml.DecodeFile(configFile, &c); err != nil {
		log.Fatalln("Failed to open file", err)
		return
	}

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

	//log.Println("File loaded for " + c.ClientId)
}

// TODO: support -c property?
func findConfigFile() string {
	// Prepare list of directories
	user, err := user.Current()
	if err != nil {
		// TODO: don't fail here, just skip locations that require the user.
		log.Fatal(err)
	}

	wd, _ := os.Getwd()

	paths := []string{
		wd, // current working directory
		"/etc/sparkcli",
		user.HomeDir, // users' home directory
	}

	for _, basepath := range paths {
		path := basepath + string(os.PathSeparator) + "sparkcli.toml"
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return "sparkcli.toml"
}

func (c Configuration) Save() {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(c); err != nil {
		log.Fatalln("Failed to encode config", err)
	}
	f, err := os.Create(configFile)
	if err != nil {
		log.Fatalln("Failed to create file", err)
		return
	}

	w := bufio.NewWriter(f)
	buf.WriteTo(w)
	w.Flush()
}

func (c Configuration) checkClientConfig() error {
	if c.ClientId == "" {
		return errors.New("ClientId not configured")
	}
	if c.ClientSecret == "" {
		return errors.New("ClientSecret not configured")
	}
	if c.AuthCode == "" {
		c.PrintAuthUrl()
		return errors.New("AuthCode not configured")
	}
	return nil
}

func (c Configuration) checkAccessToken() bool {
	if c.AccessToken == "" {
		return false
	} else {
		return true
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
