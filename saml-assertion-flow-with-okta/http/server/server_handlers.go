package server

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/emanor-okta/saml-assertion-flow-with-okta/http/client"
	config "github.com/emanor-okta/saml-assertion-flow-with-okta/utils"
	util "github.com/emanor-okta/saml-assertion-flow-with-okta/utils"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	// logger = log.New(log.Writer(), "", log.Lshortfile)
}

func RootHandler(res http.ResponseWriter, req *http.Request) {
	// data := []byte("Returned Root")
	// res.WriteHeader(200)
	// wrote, _ := res.Write(data)
	// fmt.Printf("Wrote %v bytes\n", wrote)
	tpl.ExecuteTemplate(res, "home.gohtml", config.GetConfiguration())
}

func HandleSamlResponse(res http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	samlResp := req.PostForm.Get("SAMLResponse")
	if samlResp != "" {
		client.Flowstate.SamlResp = samlResp
		client.Flowstate.SamlRespD = util.FormatXML(util.DecodeB64(samlResp))
		ass := util.GetAssertionFromOktaResponse(samlResp)
		fmt.Printf("SAML Resp:\n%s\n", ass)
		client.Flowstate.SamlAssertion = ass
		client.Flowstate.SamlAssertionD = util.FormatXML(util.URLDecodeB64(ass))
		GetTokens(res, req)
		return
	} else {
		fmt.Println("SAMLResponse Not found..")
		for k, v := range req.PostForm {
			fmt.Printf("%s - %s\n", k, v)
		}
	}
	res.Write([]byte("No SAML Response"))
}

/*
 * /gettokens - calls the /token endpoint on Okta Org, presents the call history
 */
func GetTokens(res http.ResponseWriter, req *http.Request) {
	client.GetTokens(client.Flowstate.SamlAssertion)
	// client.Kclogin.SamlReqURL = util.SAML_REQUEST_URL
	client.Flowstate.TokenURL = util.TOKEN_EP
	values := struct {
		client.FlowState
		EMBED_LINK string
	}{
		client.Flowstate,
		config.GetConfiguration().EMBED_LINK,
	}
	tpl.ExecuteTemplate(res, "tokens.gohtml", values)
}

/*
 * /config - configure Okta
 */
func ConfigHandler(res http.ResponseWriter, req *http.Request) {
	c := config.GetConfiguration()
	req.ParseForm()
	if len(req.PostForm) > 0 {
		for k, v := range req.PostForm {
			fmt.Printf("%s - %s\n", k, v)
			switch k {
			case "clientId":
				c.CLIENT_ID = v[0]
			case "clientSecret":
				c.CLIENT_SECRET = v[0]
			case "tokenURL":
				c.TOKEN_EP = v[0]
			case "samlEmbedLink":
				c.EMBED_LINK = v[0]
			}
		}
		config.SaveConfiguration(c)
	}

	tpl.ExecuteTemplate(res, "config.gohtml", c)
}
