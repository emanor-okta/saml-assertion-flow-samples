package utils

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

const RELAY_STATE string = "%2Fapp%2FUserHome"
const SAML_REQUEST string = `<?xml version="1.0" encoding="UTF-8"?><saml2p:AuthnRequest AssertionConsumerServiceURL="%s" Destination="%s" ForceAuthn="false" ID="id406023769707305891283735" IssueInstant="%s" Version="2.0" xmlns:saml2p="urn:oasis:names:tc:SAML:2.0:protocol"><saml2:Issuer xmlns:saml2="urn:oasis:names:tc:SAML:2.0:assertion">%s</saml2:Issuer><saml2p:NameIDPolicy Format="urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified"/></saml2p:AuthnRequest>`

var SAML_REQUEST_URL string
var SAML_ACS_URL string
var SAML_ISSUER string
var TOKEN_EP string
var CLIENT_ID string
var CLIENT_SECRET string

var c Configuration

// var configured bool

type Kc_idp struct {
	SAML_REQUEST_URL string
	SAML_ACS_URL     string
	SAML_ISSUER      string
}

type Okta_app struct {
	CLIENT_ID     string
	CLIENT_SECRET string
	TOKEN_EP      string
}

type Configuration struct {
	Kc_idp
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
	SAML_REQUEST_URL = c.SAML_REQUEST_URL
	SAML_ACS_URL = c.SAML_ACS_URL
	SAML_ISSUER = c.SAML_ISSUER
	TOKEN_EP = c.TOKEN_EP
	CLIENT_ID = c.CLIENT_ID
	CLIENT_SECRET = c.CLIENT_SECRET
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
	if SAML_ACS_URL != "" && SAML_ISSUER != "" && SAML_REQUEST_URL != "" &&
		TOKEN_EP != "" && CLIENT_ID != "" && CLIENT_SECRET != "" {
		c.Configured = true
	} else {
		c.Configured = false
	}
}
