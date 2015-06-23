package store

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// UninstallingPreview is the preview for a currently-uninstalling package. It's
// a placeholder for an actual progress indicator, which isn't yet available for
// this scope.
type UninstallingPreview struct {
	snap webdm.Package
}

// NewUninstallingPreview creates a new UninstallingPreview.
//
// Parameters:
// snap: Snap to be represented by this preview.
//
// Returns:
// - Pointer to new UninstallingPreview (nil if error)
// - Error (nil if none)
func NewUninstallingPreview(snap webdm.Package) (*UninstallingPreview, error) {
	if snap.NotInstalled() {
		return nil, fmt.Errorf("Snap isn't installed")
	}

	return &UninstallingPreview{snap: snap}, nil
}

// Generate pushes the preview widgets onto a WidgetReceiver.
//
// Parameters:
// receiver: Implementation of the WidgetReceiver interface.
//
// Returns:
// - Error (nil if none)
func (preview UninstallingPreview) Generate(receiver WidgetReceiver) error {
	if preview.snap.NotInstalled() {
		return fmt.Errorf("Snap isn't installed")
	}

	receiver.PushWidgets(preview.textWidget())
	receiver.PushWidgets(preview.actionsWidget())

	return nil
}

// textWidget is used to create a text widget for the uninstalling progress.
//
// Returns:
// - Text preview widget for the progress.
func (preview UninstallingPreview) textWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("progress", "text")

	// Make sure we show an error if necessary.
	if preview.snap.Message == "" {
		widget.AddAttributeValue("title", "Pretend progress was here")
		widget.AddAttributeValue("text", "Progress is still a work in, well, progress. Click this handy button to manually refresh.")
	} else {
		widget.AddAttributeValue("title", "Unable to uninstall")
		widget.AddAttributeValue("text", preview.snap.Message)
	}

	return widget
}

// actionWidget is used to create an action widget to refresh the progress.
//
// Returns:
// - Action preview widget for the progress
func (preview UninstallingPreview) actionsWidget() scopes.PreviewWidget {
	action := make(map[string]interface{})

	// Make sure we show an error if it happens.
	if preview.snap.Message == "" {
		action["label"] = "Refresh"
		action["id"] = ActionRefreshUninstalling
	} else {
		action["label"] = "Okay"
		action["id"] = ActionOk
	}

	widget := scopes.NewPreviewWidget("actions", "actions")
	widget.AddAttributeValue("actions", []interface{}{action})

	return widget
}
