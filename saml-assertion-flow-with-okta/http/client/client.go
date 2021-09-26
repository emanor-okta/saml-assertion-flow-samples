package client

import (
	// "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	// "time"

	"github.com/emanor-okta/saml-assertion-flow-with-okta/utils"
	config "github.com/emanor-okta/saml-assertion-flow-with-okta/utils"
)

var client *http.Client
var Flowstate FlowState
var ClientLogger *log.Logger

type FlowState struct {
	Url            string
	Username       string
	Password       string
	SamlResp       string
	SamlRespD      string
	SamlAssertion  string
	SamlAssertionD string
	Tokens         string
	IdToken        string
	AccessToken    string
	EmbedLink      string
	TokenURL       string
	BasicAuth      string
	Configured     bool
	Error          string
}

func init() {
	client = &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	jar, _ := cookiejar.New(nil)
	client.Jar = jar
	ClientLogger = log.New(log.Writer(), "", log.Ldate|log.Ltime|log.Lshortfile)
	// ClientLogger.Println("\n\nTESTing")
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
	Flowstate.Tokens = string(toks)
	Flowstate.AccessToken = utils.FormatJSON(utils.RawDecodeB64(strings.Split(tokens.Access_token, ".")[1]))
	Flowstate.IdToken = utils.FormatJSON(utils.RawDecodeB64(strings.Split(tokens.Id_token, ".")[1]))
	Flowstate.BasicAuth = req.Header.Get("Authorization")
	Flowstate.EmbedLink = config.EMBED_LINK
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
