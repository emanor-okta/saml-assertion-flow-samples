package main

import (
	"fmt"

	"github.com/emanor-okta/saml-assertion-flow-sample/http/server"
)

func main() {
	fmt.Println("starting...")
	server.StartServer()
}
