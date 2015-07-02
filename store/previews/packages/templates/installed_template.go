package templates

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/actions"
	"launchpad.net/unity-scope-snappy/store/previews/humanize"
	"launchpad.net/unity-scope-snappy/webdm"
)

// InstalledTemplate is a preview template for an installed package.
type InstalledTemplate struct {
	*GenericTemplate
}

// NewInstalledTemplate creates a new InstalledTemplate.
//
// Parameters:
// snap: Snap to be represented by this template.
//
// Returns:
// - Pointer to new InstalledTemplate (nil if error)
// - Error (nil if none)
func NewInstalledTemplate(snap webdm.Package) (*InstalledTemplate, error) {
	if !snap.Installed() {
		return nil, fmt.Errorf("Snap is not installed")
	}

	template := new(InstalledTemplate)
	template.GenericTemplate = NewGenericTemplate(snap)

	return template, nil
}

// actionsWidget is used to create an actions widget to uninstall/open the snap.
//
// Returns:
// - Action preview widget for the snap.
func (preview InstalledTemplate) ActionsWidget() scopes.PreviewWidget {
	widget := preview.GenericTemplate.ActionsWidget()

	openAction := make(map[string]interface{})
	openAction["id"] = actions.ActionOpen
	openAction["label"] = "Open"

	uninstallAction := make(map[string]interface{})
	uninstallAction["id"] = actions.ActionUninstall
	uninstallAction["label"] = "Uninstall"

	widget.AddAttributeValue("actions", []interface{}{openAction, uninstallAction})

	return widget
}

// updatesWidget is used to create a table widget holding snap information.
//
// Returns:
// - Table widget for the snap.
func (preview InstalledTemplate) UpdatesWidget() scopes.PreviewWidget {
	widget := preview.GenericTemplate.UpdatesWidget()

	value, ok := widget["values"]
	if ok {
		rows := value.([]interface{})
		if rows != nil {
			sizeRow := []string{"Size", humanize.Bytes(preview.snap.InstalledSize)}
			rows = append(rows, sizeRow)

			widget.AddAttributeValue("values", rows)
		}
	}

	return widget
}
