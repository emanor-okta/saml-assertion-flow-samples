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
	ClientLogger.Println("\n\nTESTing")
}

func SendSAMLRequest() *KcLogin {
	now := time.Now()
	samReq := fmt.Sprintf(config.SAML_REQUEST, config.SAML_ACS_URL, config.SAML_REQUEST_URL,
		now.Format("1970-01-01T00:00:00.000Z"), config.SAML_ISSUER)
	samReqEnc := base64.StdEncoding.EncodeToString([]byte(samReq))
	fmt.Println(samReq)
	fmt.Println(samReqEnc)
	resp, err := client.PostForm(config.SAML_REQUEST_URL, url.Values{
		"SAMLRequest": {samReqEnc},
		"RelayState":  {config.RELAY_STATE},
		// "LoginHint":   {"key1@cloak.com"},
	})
	if err != nil {
		log.Fatal(err)
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
	}

	// hack since ParseKeycloakResponse will return new instance, dont want to loose SAML Request
	samlReq := kcLogin.SamlReq
	samlReqD := kcLogin.SamlReqD
	*kcLogin = ParseKeyCloakResponse(resp.Body)
	kcLogin.SamlReq = samlReq
	kcLogin.SamlReqD = samlReqD
}

func GetTokens(a string) {
	// fmt.Printf("\n\nAssertion\n%v\n\n", a)
	v := url.Values{
		"grant_type": {"urn:ietf:params:oauth:grant-type:saml2-bearer"},
		"scope":      {"openid profile"},
		"assertion":  {a},
	}
	req, err := http.NewRequest("POST", config.TOKEN_EP, strings.NewReader(v.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(config.CLIENT_ID, config.CLIENT_SECRET)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	tokens := struct {
		Token_type   string `json:"token_type"`
		Expires_in   int    `json:"expires_in"`
		Access_token string `json:"access_token"`
		Scope        string `json:"scope"`
		Id_token     string `json:"id_token"`
	}{}
	// tokens := Tokens{}
	fmt.Println("RESPONSE:")
	fmt.Println(resp)
	res, _ := io.ReadAll(resp.Body)
	json.Unmarshal(res, &tokens)
	toks, _ := json.MarshalIndent(tokens, "", "  ")
	Kclogin.Tokens = string(toks)
	Kclogin.AccessToken = utils.FormatJSON(utils.RawDecodeB64(strings.Split(tokens.Access_token, ".")[1]))
	Kclogin.IdToken = utils.FormatJSON(utils.RawDecodeB64(strings.Split(tokens.Id_token, ".")[1]))
	Kclogin.BasicAuth = req.Header.Get("Authorization")
	fmt.Println(string(res))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	log.Println("\n\nIn REdirect...")
	log.Println()
	log.Println(*req)
	log.Println()
	log.Println(*req.Response)
	log.Println()

	return nil
}
