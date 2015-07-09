package templates

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// InstallingTemplate is a preview template for a package that is currently
// being installed. It's based upon the StoreTemplate.
type InstallingTemplate struct {
	*StoreTemplate
	objectPath dbus.ObjectPath
}

// NewInstallingTemplate creates a new InstallingTemplate.
//
// Parameters:
// snap: Snap to be represented by this template.
// objectPath: DBus object path upon which progress updates will be provided.
//
// Returns:
// - Pointer to new InstallingTemplate (nil if error)
// - Error (nil if none)
func NewInstallingTemplate(snap webdm.Package, objectPath dbus.ObjectPath) (*InstallingTemplate, error) {
	if snap.Uninstalling() {
		return nil, fmt.Errorf("Snap is currently being uninstalled")
	}

	if !objectPath.IsValid() {
		return nil, fmt.Errorf(`Invalid object path: "%s"`, objectPath)
	}

	template := &InstallingTemplate{objectPath: objectPath}

	var err error
	template.StoreTemplate, err = NewStoreTemplate(snap)
	if err != nil {
		return nil, fmt.Errorf("Unable to create store template: %s", err)
	}

	return template, nil
}

// ActionsWidget is used to create a progress widget where the store actions
// were.
//
// Returns:
// - Progress preview widget for the snap.
func (preview InstallingTemplate) ActionsWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("install", "progress")

	source := make(map[string]interface{})
	source["dbus-name"] = "com.canonical.applications.WebdmPackageManager"
	source["dbus-object"] = preview.objectPath

	widget.AddAttributeValue("source", source)

	return widget
}
