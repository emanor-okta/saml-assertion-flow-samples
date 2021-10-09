package client

import (
	"io"

	"golang.org/x/net/html"
)

type KcLogin struct {
	Url            string
	Username       string
	Password       string
	SamlResp       string
	SamlRespD      string
	SamlAssertion  string
	SamlAssertionD string
	SamlReq        string
	SamlReqD       string
	Tokens         string
	IdToken        string
	AccessToken    string
	SamlReqURL     string
	TokenURL       string
	BasicAuth      string
	Configured     bool
	Error          string
}

func ParseKeyCloakResponse(r io.Reader) KcLogin {
	kcLogin := KcLogin{}
	kcLogin.Configured = true

	doc, err := html.Parse(r)
	if err != nil {
		ClientLogger.Printf("ERROR: %s\n", err)
		kcLogin.Error = err.Error()
		return kcLogin
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "form" {
			// login form returned
			found := false
			for _, a := range n.Attr {
				if a.Key == "id" {
					if a.Val == "kc-form-login" {
						found = true
						continue
					}
				} else if a.Key == "action" {
					if found {
						kcLogin.Url = a.Val
						break
					}
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "input" {
			// SAML Response Returned (as POST Message)
			found := false
			for _, a := range n.Attr {
				if a.Key == "name" {
					if a.Val == "SAMLResponse" {
						found = true
						continue
					}
				} else if a.Key == "value" {
					if found {
						kcLogin.SamlResp = a.Val
						break
					}
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "span" {
			for _, a := range n.Attr {
				// Check for login error
				if a.Key == "id" {
					if a.Val == "input-error" {
						for c := n.FirstChild; c != nil; c = c.NextSibling {
							if c.Type == html.TextNode {
								ClientLogger.Printf("ERROR: %s\n", c.Data)
								kcLogin.Error = c.Data
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil && kcLogin.Error == ""; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return kcLogin
}
