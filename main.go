package main

import (
	"flag"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/scope"
	"log"
)

// main is the entry point of the scope.
//
// Command-line parameters:
// - webdm=address[:port] on which WebDM is listening
func main() {
	webdmAddressParameter := flag.String("webdm", "", "WebDM address[:port]")
	flag.Parse()

	scope, err := scope.NewScope(*webdmAddressParameter)
	if err != nil {
		log.Printf("unity-scope-snappy: Unable to create scope: %s", err)
		return
	}

	err = scopes.Run(scope)
	if err != nil {
		log.Printf("unity-scope-snappy: Unable to run scope: %s", err)
	}
}
