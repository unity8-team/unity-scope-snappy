package templates

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// InstalledTemplate is a preview template for a package that is currently being
// uninstalled. It's based upon the InstalledTemplate.
type UninstallingTemplate struct {
	*InstalledTemplate
	objectPath dbus.ObjectPath
}

// NewUninstallingTemplate creates a new UninstallingTemplate.
//
// Parameters:
// snap: Snap to be represented by this template.
// objectPath: DBus object path upon which progress updates will be provided.
//
// Returns:
// - Pointer to new UninstallingTemplate (nil if error)
// - Error (nil if none)
func NewUninstallingTemplate(snap webdm.Package, objectPath dbus.ObjectPath) (*UninstallingTemplate, error) {
	if snap.Installing() {
		return nil, fmt.Errorf("Snap is currently being installed")
	}

	if !objectPath.IsValid() {
		return nil, fmt.Errorf(`Invalid object path: "%s"`, objectPath)
	}

	template := &UninstallingTemplate{objectPath: objectPath}

	var err error
	template.InstalledTemplate, err = NewInstalledTemplate(snap)
	if err != nil {
		return nil, fmt.Errorf("Unable to create installed template: %s", err)
	}

	return template, nil
}

// ActionWidget is used to create a progress widget where the store actions
// were.
//
// Returns:
// - Progress preview widget for the snap.
func (preview UninstallingTemplate) ActionsWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("uninstall", "progress")

	source := make(map[string]interface{})
	source["dbus-name"] = "com.canonical.applications.WebdmPackageManager"
	source["dbus-object"] = preview.objectPath

	widget.AddAttributeValue("source", source)

	return widget
}
