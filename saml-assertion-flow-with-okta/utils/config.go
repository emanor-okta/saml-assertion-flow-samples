package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

var TOKEN_EP string
var CLIENT_ID string
var CLIENT_SECRET string
var EMBED_LINK string

var c Configuration

type Okta_app struct {
	CLIENT_ID     string
	CLIENT_SECRET string
	TOKEN_EP      string
	EMBED_LINK    string
}

type Configuration struct {
	Okta_app
	Configured bool
}

func init() {
	loadConfig(&c)
	// SaveConfiguration(&c)
	setVars(&c)
	setIsConfigured()
}

func GetConfiguration() *Configuration {
	return &c
}

func setVars(c *Configuration) {
	TOKEN_EP = c.TOKEN_EP
	CLIENT_ID = c.CLIENT_ID
	CLIENT_SECRET = c.CLIENT_SECRET
	EMBED_LINK = c.EMBED_LINK
}

func SaveConfiguration(c *Configuration) {
	setVars(c)
	setIsConfigured()
	bytes, err := yaml.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("config.yaml", bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func loadConfig(c *Configuration) {
	buf, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("No Configuration file exists yet: %v\n", err)
		buf = []byte{}
	}

	err = yaml.Unmarshal(buf, c)
	if err != nil {
		log.Fatal(err)
	}
}

func setIsConfigured() {
	if EMBED_LINK != "" && TOKEN_EP != "" && CLIENT_ID != "" && CLIENT_SECRET != "" {
		c.Configured = true
	} else {
		c.Configured = false
	}
}
