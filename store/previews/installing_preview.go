package previews

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/actions"
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
	"launchpad.net/unity-scope-snappy/webdm"
)

// InstallingPreview is the preview for a currently-installing package. It's
// a placeholder for an actual progress indicator, which isn't yet available for
// this scope.
type InstallingPreview struct {
	snap webdm.Package
}

// NewInstallingPreview creates a new InstallingPreview.
//
// Parameters:
// snap: Snap to be represented by this preview.
//
// Returns:
// - Pointer to new InstallingPreview (nil if error)
// - Error (nil if none)
func NewInstallingPreview(snap webdm.Package) (*InstallingPreview, error) {
	if snap.Installed() {
		return nil, fmt.Errorf("Snap is already installed")
	}

	return &InstallingPreview{snap: snap}, nil
}

// Generate pushes the preview widgets onto a WidgetReceiver.
//
// Parameters:
// receiver: Implementation of the interfaces.WidgetReceiver interface.
//
// Returns:
// - Error (nil if none)
func (preview InstallingPreview) Generate(receiver interfaces.WidgetReceiver) error {
	if preview.snap.Installed() {
		return fmt.Errorf("Snap is already installed")
	}

	receiver.PushWidgets(preview.textWidget())
	receiver.PushWidgets(preview.actionsWidget())

	return nil
}

// textWidget is used to create a text widget for the installing progress.
//
// Returns:
// - Text preview widget for the progress.
func (preview InstallingPreview) textWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("progress", "text")

	// Make sure we show an error if necessary.
	if preview.snap.Message == "" {
		widget.AddAttributeValue("title", "Pretend progress was here")
		widget.AddAttributeValue("text", "Progress is still a work in, well, progress. Click this handy button to manually refresh.")
	} else {
		widget.AddAttributeValue("title", "Unable to install")
		widget.AddAttributeValue("text", preview.snap.Message)
	}

	return widget
}

// actionWidget is used to create an action widget to refresh the progress.
//
// Returns:
// - Action preview widget for the progress.
func (preview InstallingPreview) actionsWidget() scopes.PreviewWidget {
	action := make(map[string]interface{})

	// Make sure we show an error if it happens.
	if preview.snap.Message == "" {
		action["label"] = "Refresh"
		action["id"] = actions.ActionRefreshInstalling
	} else {
		action["label"] = "Okay"
		action["id"] = actions.ActionOk
	}

	widget := scopes.NewPreviewWidget("actions", "actions")
	widget.AddAttributeValue("actions", []interface{}{action})

	return widget
}
