package main

import (
	"flag"
	"fmt"
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"log"
)

type SnappyScope struct {
	webdmClient *webdm.Client
}

func (scope *SnappyScope) SetScopeBase(base *scopes.ScopeBase) {
	// Do nothing
}

const template = `{
	"schema-version": 1,
	"template": {
		"category-layout": "grid",
		"card-size": "small"
	},
	"components": {
		"title": "title",
		"art" : {
			"field": "art"
		},
		"subtitle": "subtitle"
        }
}`

func (scope SnappyScope) Search(query *scopes.CannedQuery, metadata *scopes.SearchMetadata, reply *scopes.SearchReply, cancelled <-chan bool) error {
	packages, err := scope.webdmClient.GetStorePackages()
	if err != nil {
		errorString := fmt.Sprintf("unity-scope-snappy: Unable to retrieve store packages: %s",
			err)

		// Log to stderr as well, since nothing seems to do anything with
		// returned error (yet)
		log.Println(errorString)
		return fmt.Errorf(errorString)
	}

	category := reply.RegisterCategory("store_packages", "Store Packages", "", template)
	for _, thisPackage := range packages {
		result := scopes.NewCategorisedResult(category)

		result.SetTitle(thisPackage.Name)
		result.Set("subtitle", thisPackage.Description)
		result.SetURI("snappy:" + thisPackage.Id)
		result.SetArt(thisPackage.IconUrl)

		if reply.Push(result) != nil {
			// If the push fails, the query was cancelled. No need to continue.
			return nil
		}
	}

	return nil
}

func (scope SnappyScope) Preview(result *scopes.Result, metadata *scopes.ActionMetadata, reply *scopes.PreviewReply, cancelled <-chan bool) error {
	layout1Column := scopes.NewColumnLayout(1)
	layout2Column := scopes.NewColumnLayout(2)
	layout3Column := scopes.NewColumnLayout(3)

	layout1Column.AddColumn("image", "header", "summary")

	layout2Column.AddColumn("image")
	layout2Column.AddColumn("header", "summary")

	layout3Column.AddColumn("image")
	layout3Column.AddColumn("header", "summary")
	layout3Column.AddColumn("")

	reply.RegisterLayout(layout1Column, layout2Column, layout3Column)

	header := scopes.NewPreviewWidget("header", "header")
	header.AddAttributeMapping("title", "title")
	header.AddAttributeMapping("subtitle", "subtitle")

	image := scopes.NewPreviewWidget("image", "image")
	image.AddAttributeMapping("source", "art")

	description := scopes.NewPreviewWidget("summary", "text")
	description.AddAttributeMapping("text", "description")

	reply.PushWidgets(image, header, description)

	return nil
}

func main() {
	webdmAddressParameter := flag.String("webdm", "127.0.0.1:4200", "WebDM address[:port]")
	flag.Parse()

	scope := &SnappyScope{webdmClient: webdm.NewClient()}
	scope.webdmClient.BaseUrl.Host = *webdmAddressParameter

	err := scopes.Run(scope)
	if err != nil {
		log.Printf("unity-scope-snappy: Unable to run scope: %s", err)
	}
}
