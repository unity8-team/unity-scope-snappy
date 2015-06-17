package scope

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// InstalledPackagePreviewTemplate is a preview template for an installed
// package.
type InstalledPackagePreviewTemplate struct {
	*GenericPackagePreviewTemplate
}

// NewInstalledPackagePreviewTemplate creates a new
// InstalledPackagePreviewTemplate.
//
// Parameters:
// snap: Snap to be represented by this template.
//
// Returns:
// - Pointer to new InstalledPackagePreviewTemplate (nil if error)
// - Error (nil if none)
func NewInstalledPackagePreviewTemplate(snap webdm.Package) (*InstalledPackagePreviewTemplate, error) {
	if !snap.Installed() {
		return nil, fmt.Errorf("Snap is not installed")
	}

	template := new(InstalledPackagePreviewTemplate)
	template.GenericPackagePreviewTemplate = NewGenericPackagePreviewTemplate(snap)

	return template, nil
}

// actionsWidget is used to create an actions widget to uninstall/open the snap.
//
// Returns:
// - Action preview widget for the snap.
func (preview InstalledPackagePreviewTemplate) actionsWidget() scopes.PreviewWidget {
	widget := preview.GenericPackagePreviewTemplate.actionsWidget()

	openAction := make(map[string]interface{})
	openAction["id"] = ActionOpen
	openAction["label"] = "Open"

	uninstallAction := make(map[string]interface{})
	uninstallAction["id"] = ActionUninstall
	uninstallAction["label"] = "Uninstall"

	widget.AddAttributeValue("actions", []interface{}{openAction, uninstallAction})

	return widget
}

// updatesWidget is used to create a table widget holding snap information.
//
// Returns:
// - Table widget for the snap.
func (preview InstalledPackagePreviewTemplate) updatesWidget() scopes.PreviewWidget {
	widget := preview.GenericPackagePreviewTemplate.updatesWidget()

	value, ok := widget["values"]
	if ok {
		rows := value.([]interface{})
		if rows != nil {
			sizeRow := []string{"Size", humanizeBytes(preview.snap.InstalledSize)}
			rows = append(rows, sizeRow)

			widget.AddAttributeValue("values", rows)
		}
	}

	return widget
}
