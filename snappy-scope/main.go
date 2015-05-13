package main

import (
    "fmt"
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

type SnappyScope struct {}

func (scope SnappyScope) SetScopeBase(base *scopes.ScopeBase) {
	// Do nothing
}

const template = `{
	"schema-version": 1,
	"template": {
		"category-layout": "grid",
		"card-size": "medium"
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
	categoryTitle := "cat title"
	category := reply.RegisterCategory("current", categoryTitle, "", template)
	result := scopes.NewCategorisedResult(category)
	result.SetTitle(client.GetSomethingUseful())
	result.Set("subtitle", "fake subtitle")
	result.Set("description", "A description of the result")
	result.SetURI("http://fake")
	result.SetArt("file:/usr/share/icons/Humanity/apps/48/system-software-install.svg")

	if reply.Push(result) != nil {
		// If the push fails, the query was cancelled. No need to continue.
		return nil
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
	err := scopes.Run(&SnappyScope{})
	if err != nil {
		fmt.Println(err)
	}
}
