/* Copyright (C) 2015 Canonical Ltd.
 *
 * This file is part of unity-scope-snappy.
 *
 * unity-scope-snappy is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option) any
 * later version.
 *
 * unity-scope-snappy is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * unity-scope-snappy. If not, see <http://www.gnu.org/licenses/>.
 */

package previews

import (
	"fmt"
	"github.com/snapcore/snapd/client"
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/actions"
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
)

// ConfirmUninstallPreview is a PreviewGenerator meant to have the user
// confirm a request to uninstall a package.
type ConfirmUninstallPreview struct {
	snap client.Snap
}

// NewConfirmUninstallPreview creates a new ConfirmUninstallPreview.
//
// Parameters:
// snap: Package which we're being asked to uninstall.
func NewConfirmUninstallPreview(snap client.Snap) *ConfirmUninstallPreview {
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
