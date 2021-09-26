package main

import (
	"fmt"

	"github.com/emanor-okta/saml-assertion-flow-sample/http/server"
)

func main() {
	fmt.Println("starting...")

	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// kcLogin := client.SendSAMLRequest()
	// kcLogin.Username = "key1@cloak.com"
	// kcLogin.Password = "P@ssw0rd"
	// assertion := client.LoginKC(kcLogin)
	// client.GetTokens(assertion)

	server.StartServer()
	// testConfig()
	// wg.Wait()
}

// func testConfig() {
// 	type Kc struct {
// 		Url string
// 	}
// 	type Okta_app struct {
// 		Url string
// 	}
// 	type configuration struct {
// 		Version string
// 		Kc
// 		Okta_app
// 	}

// 	buf, err := ioutil.ReadFile("config.yaml")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	c := &configuration{}
// 	err = yaml.Unmarshal(buf, c)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("%+v", *c)
// 	bytes, err := yaml.Marshal(c)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Marshall:\n%s\n", string(bytes))
// }
