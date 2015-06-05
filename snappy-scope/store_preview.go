package main

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// StorePreview is the preview for a package that isn't installed.
type StorePreview struct {
	snap webdm.Package
}

// NewStorePreview creates a new StorePreview.
//
// Parameters:
// snap: Snap to be represented by this preview.
//
// Returns:
// - Pointer to new StorePreview (nil if error)
// - Error (nil if none)
func NewStorePreview(snap webdm.Package) (*StorePreview, error) {
	if snap.Installed() {
		return nil, fmt.Errorf("Snap is installed")
	}

	return &StorePreview{snap: snap}, nil
}

// Generate pushes the preview widgets onto a WidgetReceiver.
//
// Parameters:
// receiver: Implementation of the WidgetReceiver interface.
//
// Returns:
// - Error (nil if none)
func (preview StorePreview) Generate(receiver WidgetReceiver) error {
	if preview.snap.Installed() {
		return fmt.Errorf("Snap is installed")
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
func (preview StorePreview) headerWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("header", "header")

	widget.AddAttributeMapping("title", "title")
	widget.AddAttributeMapping("subtitle", "subtitle")
	widget.AddAttributeMapping("mascot", "art")

	priceAttribute := make(map[string]interface{})
	priceAttribute["value"] = "FREE" // All the snaps are currently free
	widget.AddAttributeValue("attributes", []interface{}{priceAttribute})

	return widget
}

// actionWidget is used to create an action widget to install the snap.
//
// Returns:
// - Action preview widget for the snap.
func (preview StorePreview) actionWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("actions", "actions")

	installAction := make(map[string]interface{})
	installAction["id"] = ActionInstall
	installAction["label"] = "Install"

	widget.AddAttributeValue("actions", []interface{}{installAction})

	return widget
}

// infoWidget is used to create a text widget holding snap information.
//
// Returns:
// - Text preview widget for the snap.
func (preview StorePreview) infoWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("summary", "text")

	widget.AddAttributeValue("title", "Info")
	widget.AddAttributeValue("text", preview.snap.Description)

	return widget
}

// updatesWidget is used to create a table widget holding snap information.
//
// Returns:
// - Table widget for the snap.
func (preview StorePreview) updatesWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("updates_table", "table")
	widget.AddAttributeValue("title", "Updates")

	versionRow := []string{"Version number", preview.snap.Version}

	widget.AddAttributeValue("values", []interface{}{versionRow})

	return widget
}
