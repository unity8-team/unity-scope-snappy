package scope

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// StorePackagePreviewTemplate is a preview template for an in-store package
// (i.e. not installed).
type StorePackagePreviewTemplate struct {
	*GenericPackagePreviewTemplate
}

// NewStorePackagePreviewTemplate creates a new StorePackagePreviewTemplate.
//
// Parameters:
// snap: Snap to be represented by this template.
//
// Returns:
// - Pointer to new StorePackagePreviewTemplate (nil if error)
// - Error (nil if none)
func NewStorePackagePreviewTemplate(snap webdm.Package) (*StorePackagePreviewTemplate, error) {
	if snap.Installed() {
		return nil, fmt.Errorf("Snap is installed")
	}

	template := new(StorePackagePreviewTemplate)
	template.GenericPackagePreviewTemplate = NewGenericPackagePreviewTemplate(snap)

	return template, nil
}

// headerWidget is used to create a header widget for the snap, including the
// price of the package.
//
// Returns:
// - Header preview widget for the snap.
func (preview StorePackagePreviewTemplate) headerWidget() scopes.PreviewWidget {
	widget := preview.GenericPackagePreviewTemplate.headerWidget()

	// WebDM doesn't provide any information about the cost of apps... so all
	// the snaps are free!
	priceAttribute := make(map[string]interface{})
	priceAttribute["value"] = "FREE"
	widget.AddAttributeValue("attributes", []interface{}{priceAttribute})

	return widget
}

// actionWidget is used to create an action widget to install the snap.
//
// Returns:
// - Action preview widget for the snap.
func (preview StorePackagePreviewTemplate) actionsWidget() scopes.PreviewWidget {
	widget := preview.GenericPackagePreviewTemplate.actionsWidget()

	installAction := make(map[string]interface{})
	installAction["id"] = ActionInstall
	installAction["label"] = "Install"

	widget.AddAttributeValue("actions", []interface{}{installAction})

	return widget
}

// updatesWidget is used to create a table widget holding snap information.
//
// Returns:
// - Table widget for the snap.
func (preview StorePackagePreviewTemplate) updatesWidget() scopes.PreviewWidget {
	widget := preview.GenericPackagePreviewTemplate.updatesWidget()

	value, ok := widget["values"]
	if ok {
		rows := value.([]interface{})
		if rows != nil {
			sizeRow := []string{"Size", humanizeBytes(preview.snap.DownloadSize)}
			rows = append(rows, sizeRow)

			widget.AddAttributeValue("values", rows)
		}
	}

	return widget
}
