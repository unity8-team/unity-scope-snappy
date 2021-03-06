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

package templates

import (
	"fmt"

	"github.com/snapcore/snapd/client"
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/actions"
	"launchpad.net/unity-scope-snappy/store/previews/humanize"
)

// InstalledTemplate is a preview template for an installed package.
type InstalledTemplate struct {
	*GenericTemplate
	snap client.Snap
}

// NewInstalledTemplate creates a new InstalledTemplate.
//
// Parameters:
// snap: Snap to be represented by this template.
//
// Returns:
// - Pointer to new InstalledTemplate (nil if error)
// - Error (nil if none)
func NewInstalledTemplate(snap client.Snap) (*InstalledTemplate, error) {
	template := new(InstalledTemplate)
	template.GenericTemplate = NewGenericTemplate(snap)
	template.snap = snap

	return template, nil
}

// HeaderWidget is used to create a header widget for the snap, including the
// fact that it's installed or purchased.
//
// Returns:
// - Header preview widget for the snap.
func (preview InstalledTemplate) HeaderWidget() scopes.PreviewWidget {
	widget := preview.GenericTemplate.HeaderWidget()

	priceAttribute := make(map[string]interface{})
	priceAttribute["value"] = "✔ INSTALLED"
	widget.AddAttributeValue("attributes", []interface{}{priceAttribute})

	return widget
}

// ActionsWidget is used to create an actions widget to uninstall/open the snap.
//
// Returns:
// - Action preview widget for the snap.
func (preview InstalledTemplate) ActionsWidget() scopes.PreviewWidget {
	widget := preview.GenericTemplate.ActionsWidget()

	previewActions := make([]interface{}, 0)

	// Only show Open if we can get the URI
	if len(preview.snap.Apps) != 0 {
		uri := fmt.Sprintf("appid://%s/%s/current-user-version",
			preview.snap.Name, preview.snap.Apps[0].Name)

		openAction := make(map[string]interface{})
		openAction["id"] = actions.ActionOpen
		openAction["label"] = "Open"
		openAction["uri"] = uri
		previewActions = append(previewActions, openAction)
	}

	uninstallAction := make(map[string]interface{})
	uninstallAction["id"] = actions.ActionUninstall
	uninstallAction["label"] = "Uninstall"
	previewActions = append(previewActions, uninstallAction)

	widget.AddAttributeValue("actions", previewActions)

	return widget
}

// UpdatesWidget is used to create a table widget holding snap information.
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
