package server

import (
	"log"
	"net/http"
)

var ServerLogger *log.Logger

func StartServer() {
	http.HandleFunc("/", RootHandler)
	// http.HandleFunc("/getassertion", GetAssertionHandler)
	http.HandleFunc("/gettokens", GetTokens)
	http.HandleFunc("/config", ConfigHandler)
	// http.HandleFunc("/login", LoginHandler)

	http.HandleFunc("/samlresponse", HandleSamlResponse)

	ServerLogger = log.New(log.Writer(), "info: ", log.Ldate|log.Ltime|log.Lshortfile)
	ServerLogger.Println("Starting on 8080")

	/*go*/
	func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			ServerLogger.Println(err)
			log.Fatalf("Server startup failed: %s\n", err)
		}
	}()
}
