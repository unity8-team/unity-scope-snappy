package scope

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// GenericPackagePreviewTemplate is a PackagePreviewTemplate implementation that
// doesn't contain any conditionals depending on package information. It's meant
// to be embedded in other structs and further specialized.
type GenericPackagePreviewTemplate struct {
	snap webdm.Package
}

// NewGenericPackagePreviewTemplate creates a new GenericPackagePreviewTemplate.
func NewGenericPackagePreviewTemplate(snap webdm.Package) *GenericPackagePreviewTemplate {
	return &GenericPackagePreviewTemplate{snap: snap}
}

// headerWidget is used to create a header widget for the snap.
//
// Returns:
// - Header preview widget for the snap.
func (preview GenericPackagePreviewTemplate) headerWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("header", "header")

	widget.AddAttributeMapping("title", "title")
	widget.AddAttributeMapping("subtitle", "subtitle")
	widget.AddAttributeMapping("mascot", "art")

	return widget
}

// actionsWidget is used to create an action widget for the snap. The widget
// contains no actions.
//
// Returns:
// - Empty action preview widget for the snap.
func (preview GenericPackagePreviewTemplate) actionsWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("actions", "actions")

	return widget
}

// infoWidget is used to create a text widget holding the snap description.
//
// Returns:
// - Text preview widget for the snap.
func (preview GenericPackagePreviewTemplate) infoWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("summary", "text")

	widget.AddAttributeValue("title", "Info")
	widget.AddAttributeValue("text", preview.snap.Description)

	return widget
}

// updatesWidget is used to create a table widget holding snap version
// information.
//
// Returns:
// - Table widget for the snap.
func (preview GenericPackagePreviewTemplate) updatesWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("updates_table", "table")
	widget.AddAttributeValue("title", "Updates")

	versionRow := []string{"Version number", preview.snap.Version}

	widget.AddAttributeValue("values", []interface{}{versionRow})

	return widget
}
