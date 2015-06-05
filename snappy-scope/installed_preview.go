package main

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// InstalledPreview is the preview for an installed package.
type InstalledPreview struct {
	snap webdm.Package
}

// NewInstalledPreview creates a new InstalledPreview.
//
// Parameters:
// snap: Snap to be represented by this preview.
//
// Returns:
// - Pointer to new InstalledPreview (nil if error)
// - Error (nil if none)
func NewInstalledPreview(snap webdm.Package) (*InstalledPreview, error) {
	if !snap.Installed() {
		return nil, fmt.Errorf("Snap is not installed")
	}

	return &InstalledPreview{snap: snap}, nil
}

// Generate pushes the preview widgets onto a WidgetReceiver.
//
// Parameters:
// receiver: Implementation of the WidgetReceiver interface.
//
// Returns:
// - Error (nil if none)
func (preview InstalledPreview) Generate(receiver WidgetReceiver) error {
	if !preview.snap.Installed() {
		return fmt.Errorf("Snap is not installed")
	}

	receiver.PushWidgets(preview.headerWidget())
	receiver.PushWidgets(preview.actionWidget())
	receiver.PushWidgets(preview.infoWidget())
	receiver.PushWidgets(preview.updatesWidget())

	return nil
}

// headerWidget is used to create a header widget for the snap.
//
// Returns:
// - Header preview widget for the snap.
func (preview InstalledPreview) headerWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("header", "header")

	widget.AddAttributeMapping("title", "title")
	widget.AddAttributeMapping("subtitle", "subtitle")
	widget.AddAttributeMapping("mascot", "art")

	return widget
}

// actionWidget is used to create an action widget to uninstall/open the snap.
//
// Returns:
// - Action preview widget for the snap.
func (preview InstalledPreview) actionWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("actions", "actions")

	openAction := make(map[string]interface{})
	openAction["id"] = ActionOpen
	openAction["label"] = "Open"

	uninstallAction := make(map[string]interface{})
	uninstallAction["id"] = ActionUninstall
	uninstallAction["label"] = "Uninstall"

	widget.AddAttributeValue("actions", []interface{}{openAction, uninstallAction})

	return widget
}

// infoWidget is used to create a text widget holding snap information.
//
// Returns:
// - Text preview widget for the snap.
func (preview InstalledPreview) infoWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("summary", "text")
	widget.AddAttributeValue("title", "Info")
	widget.AddAttributeValue("text", preview.snap.Description)

	return widget
}

// updatesWidget is used to create a table widget holding snap information.
//
// Returns:
// - Table widget for the snap.
func (preview InstalledPreview) updatesWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("updates_table", "table")
	widget.AddAttributeValue("title", "Updates")

	versionRow := []string{"Version number", preview.snap.Version}

	widget.AddAttributeValue("values", []interface{}{versionRow})

	return widget
}
