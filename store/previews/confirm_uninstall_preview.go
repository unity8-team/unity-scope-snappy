package previews

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/actions"
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
	"launchpad.net/unity-scope-snappy/webdm"
)

// ConfirmUninstallPreview is a PreviewGenerator meant to have the user
// confirm a request to uninstall a package.
type ConfirmUninstallPreview struct {
	snap webdm.Package
}

// NewConfirmUninstallPreview creates a new ConfirmUninstallPreview.
//
// Parameters:
// snap: Package which we're being asked to uninstall.
func NewConfirmUninstallPreview(snap webdm.Package) *ConfirmUninstallPreview {
	return &ConfirmUninstallPreview{snap: snap}
}

// Generate pushes the template's preview widgets onto a WidgetReceiver.
//
// Parameters:
// receiver: Implementation of the WidgetReceiver interface.
//
// Returns:
// - Error (nil if none)
func (preview ConfirmUninstallPreview) Generate(receiver interfaces.WidgetReceiver) error {
	receiver.PushWidgets(preview.textWidget())
	receiver.PushWidgets(preview.actionsWidget())

	return nil
}

// textWidget is used to create a text widget for the installing progress.
//
// Returns:
// - Text preview widget for the progress.
func (preview ConfirmUninstallPreview) textWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("confirm", "text")

	widget.AddAttributeValue("text", fmt.Sprintf("Are you sure you want to uninstall %s?", preview.snap.Name))

	return widget
}

// actionsWidget is used to create an action widget to refresh the progress.
//
// Returns:
// - Action preview widget for the progress.
func (preview ConfirmUninstallPreview) actionsWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("confirmation", "actions")

	uninstallConfirmAction := make(map[string]interface{})
	uninstallConfirmAction["id"] = actions.ActionUninstallConfirm
	uninstallConfirmAction["label"] = "Uninstall"

	uninstallCancelAction := make(map[string]interface{})
	uninstallCancelAction["id"] = actions.ActionUninstallCancel
	uninstallCancelAction["label"] = "Cancel"

	widget.AddAttributeValue("actions", []interface{}{uninstallConfirmAction, uninstallCancelAction})

	return widget
}
