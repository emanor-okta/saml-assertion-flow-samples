package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/emanor-okta/saml-assertion-flow-sample/http/client"
	config "github.com/emanor-okta/saml-assertion-flow-sample/utils"
	util "github.com/emanor-okta/saml-assertion-flow-sample/utils"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func RootHandler(res http.ResponseWriter, req *http.Request) {
	tpl.ExecuteTemplate(res, "home.gohtml", config.GetConfiguration())
}

/*
 * /getassertion - get assertion from Keycloak, show login if no session exists
 */
func GetAssertionHandler(res http.ResponseWriter, req *http.Request) {
	client.Kclogin = *client.SendSAMLRequest()
	if client.Kclogin.Error != "" {
		tpl.ExecuteTemplate(res, "error.gohtml", client.Kclogin)
	} else if client.Kclogin.SamlResp == "" {
		tpl.ExecuteTemplate(res, "loginform.gohtml", client.Kclogin)
	} else {
		completeAssertionHanlder(res, req, &client.Kclogin)
	}
}

/*
 * /login - login to Keycloak with provided credentials
 */
func LoginHandler(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	client.Kclogin.Username = req.PostFormValue("username")
	client.Kclogin.Password = req.PostFormValue("password")
	client.LoginKC(&client.Kclogin)
	completeAssertionHanlder(res, req, &client.Kclogin)
}

/*
 * NO route - called from (GetAssertionHandler | LoginHandler) finish getting assertion redir to /gettokens
 */
func completeAssertionHanlder(res http.ResponseWriter, req *http.Request, kcLogin *client.KcLogin) {
	if kcLogin.Error != "" {
		ServerLogger.Printf("\nError - %s\n", kcLogin.Error)
		tpl.ExecuteTemplate(res, "error.gohtml", kcLogin)
		return
	}
	if kcLogin.SamlResp == "" {
		ServerLogger.Printf("Error - SAML Response not returned.\n%s\n", kcLogin.Error)
		tpl.ExecuteTemplate(res, "error.gohtml", kcLogin)
		return
	}
	kcLogin.SamlRespD = util.FormatXML(util.DecodeB64(kcLogin.SamlResp))
	kcLogin.SamlAssertion = util.GetAssertionFromKCResponse(kcLogin.SamlResp)
	kcLogin.SamlAssertionD = util.FormatXML(util.URLDecodeB64(kcLogin.SamlAssertion))
	http.Redirect(res, req, "/gettokens", http.StatusFound)
}

func HandleSamlResponse(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	samlResp := req.PostForm.Get("SAMLResponse")
	if samlResp != "" {
		sa := util.GetAssertionFromOktaResponse(samlResp)
		ServerLogger.Printf("SAML Resp:\n%s\n", sa)
		client.Kclogin.SamlAssertion = sa
		GetTokens(res, req)
		return
	} else {
		fmt.Println("SAMLResponse Not found..")
		for k, v := range req.PostForm {
			ServerLogger.Printf("%s - %s\n", k, v)
		}
	}
	res.Write([]byte("SAML Response"))
}

/*
 * /gettokens - calls the /token endpoint then render token template
 */
func GetTokens(res http.ResponseWriter, req *http.Request) {
	client.GetTokens(client.Kclogin.SamlAssertion)

	if client.Kclogin.Error != "" {
		tpl.ExecuteTemplate(res, "error.gohtml", client.Kclogin)
	} else {
		client.Kclogin.SamlReqURL = util.SAML_REQUEST_URL
		client.Kclogin.TokenURL = util.TOKEN_EP
		tpl.ExecuteTemplate(res, "tokens.gohtml", client.Kclogin)
	}
}

/*
 * /config - configure Keycloak and Okta
 */
func ConfigHandler(res http.ResponseWriter, req *http.Request) {
	c := config.GetConfiguration()
	req.ParseForm()
	if len(req.PostForm) > 0 {
		for k, v := range req.PostForm {
			fmt.Printf("%s - %s\n", k, v)
			switch k {
			case "samlRequestURL":
				c.SAML_REQUEST_URL = v[0]
			case "samlAcsURL":
				c.SAML_ACS_URL = v[0]
			case "samlIssuer":
				c.SAML_ISSUER = v[0]
			case "clientId":
				c.CLIENT_ID = v[0]
			case "clientSecret":
				c.CLIENT_SECRET = v[0]
			case "tokenURL":
				c.TOKEN_EP = v[0]
			}
		}
		config.SaveConfiguration(c)
	}

	tpl.ExecuteTemplate(res, "config.gohtml", c)
}
