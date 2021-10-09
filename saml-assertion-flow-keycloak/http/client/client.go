package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/emanor-okta/saml-assertion-flow-sample/utils"
	config "github.com/emanor-okta/saml-assertion-flow-sample/utils"
)

var client *http.Client
var Kclogin KcLogin
var ClientLogger *log.Logger

func init() {
	client = &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	ClientLogger = log.New(log.Writer(), "", log.Ldate|log.Ltime|log.Lshortfile)
}

func SendSAMLRequest() *KcLogin {
	now := time.Now()
	samReq := fmt.Sprintf(config.SAML_REQUEST, config.SAML_ACS_URL, config.SAML_REQUEST_URL,
		now.Format("1970-01-01T00:00:00.000Z"), config.SAML_ISSUER)
	samReqEnc := base64.StdEncoding.EncodeToString([]byte(samReq))
	ClientLogger.Println("Sending SAML Request")
	fmt.Println(samReq)
	fmt.Println(samReqEnc)

	resp, err := client.PostForm(config.SAML_REQUEST_URL, url.Values{
		"SAMLRequest": {samReqEnc},
		"RelayState":  {config.RELAY_STATE},
	})
	if err != nil {
		ClientLogger.Printf("Error Sending SAML Request:\n%s\n", err.Error())
		kcLogin := KcLogin{}
		kcLogin.Error = err.Error()
		return &kcLogin
	}
	defer resp.Body.Close()
	Kclogin = ParseKeyCloakResponse(resp.Body)
	if Kclogin.Error == "" {
		Kclogin.SamlReq = samReqEnc
		Kclogin.SamlReqD = utils.FormatXML(samReq)
	}

	ClientLogger.Printf("Kclogin:\n%v\n", Kclogin)
	return &Kclogin
}

func LoginKC(kcLogin *KcLogin) {
	resp, err := client.PostForm(kcLogin.Url, url.Values{
		"username": {kcLogin.Username},
		"password": {kcLogin.Password},
	})
	if err != nil {
		ClientLogger.Printf("Error:\n%s\n", err.Error())
		kcLogin.Error = err.Error()
		return
	}

	// since ParseKeycloakResponse will return new instance, Add SAML Request
	samlReq := kcLogin.SamlReq
	samlReqD := kcLogin.SamlReqD
	*kcLogin = ParseKeyCloakResponse(resp.Body)
	kcLogin.SamlReq = samlReq
	kcLogin.SamlReqD = samlReqD
}

func GetTokens(a string) {
	fmt.Printf("\n\nAssertion\n%v\n\n", a)
	v := url.Values{
		"grant_type": {"urn:ietf:params:oauth:grant-type:saml2-bearer"},
		"scope":      {"openid profile"},
		"assertion":  {a},
	}
	req, err := http.NewRequest("POST", config.TOKEN_EP, strings.NewReader(v.Encode()))
	if err != nil {
		ClientLogger.Println(err.Error())
		Kclogin.Error = err.Error()
		return
	}

	req.SetBasicAuth(config.CLIENT_ID, config.CLIENT_SECRET)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		ClientLogger.Println(err.Error())
		Kclogin.Error = err.Error()
		return
	}
	defer resp.Body.Close()
	tokens := struct {
		Token_type   string `json:"token_type"`
		Expires_in   int    `json:"expires_in"`
		Access_token string `json:"access_token"`
		Scope        string `json:"scope"`
		Id_token     string `json:"id_token"`
	}{}

	ClientLogger.Printf("\n/Token Response:\nStatus Code: %v\n%v\n", resp.StatusCode, resp)
	res, _ := io.ReadAll(resp.Body)
	json.Unmarshal(res, &tokens)
	toks, _ := json.MarshalIndent(tokens, "", "  ")
	if resp.StatusCode == 200 {
		Kclogin.Tokens = string(toks)
		Kclogin.AccessToken = utils.FormatJSON(utils.RawDecodeB64(strings.Split(tokens.Access_token, ".")[1]))
		Kclogin.IdToken = utils.FormatJSON(utils.RawDecodeB64(strings.Split(tokens.Id_token, ".")[1]))
	} else {
		Kclogin.Tokens = resp.Status + ": " + string(res)
	}
	Kclogin.BasicAuth = req.Header.Get("Authorization")
	ClientLogger.Println(string(res))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	log.Printf("\n\nIn Redirect\nRequest:\n%v\nResponse:\n%v\n\n", *req, *req.Response)
	return nil
}
