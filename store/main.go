package main

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/scope"
	"log"
	"os"
)

// main is the entry point of the scope.
//
// Supported environment variables:
// - WEBDM_URL: address[:port] on which WebDM is listening
func main() {
	scope, err := scope.New(os.Getenv("WEBDM_URL"))
	if err != nil {
		log.Printf("unity-scope-snappy: Unable to create scope: %s", err)
		return
	}

	err = scopes.Run(scope)
	if err != nil {
		log.Printf("unity-scope-snappy: Unable to run scope: %s", err)
	}
}