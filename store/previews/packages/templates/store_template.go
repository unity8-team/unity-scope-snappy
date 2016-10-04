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
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/actions"
	"launchpad.net/unity-scope-snappy/store/previews/humanize"
	"launchpad.net/unity-scope-snappy/webdm"
)

// StoreTemplate is a preview template for an in-store package
// (i.e. not installed).
type StoreTemplate struct {
	*GenericTemplate
}

// NewStoreTemplate creates a new StoreTemplate.
//
// Parameters:
// snap: Snap to be represented by this template.
//
// Returns:
// - Pointer to new StoreTemplate (nil if error)
// - Error (nil if none)
func NewStoreTemplate(snap webdm.Package) (*StoreTemplate, error) {
	if snap.Installed() {
		return nil, fmt.Errorf("Snap is installed")
	}

	template := new(StoreTemplate)
	template.GenericTemplate = NewGenericTemplate(snap)

	return template, nil
}

// HeaderWidget is used to create a header widget for the snap, including the
// price of the package.
//
// Returns:
// - Header preview widget for the snap.
func (preview StoreTemplate) HeaderWidget() scopes.PreviewWidget {
	widget := preview.GenericTemplate.HeaderWidget()

	// WebDM doesn't provide any information about the cost of apps... so all
	// the snaps are free!
	priceAttribute := make(map[string]interface{})
	priceAttribute["value"] = "FREE"
	widget.AddAttributeValue("attributes", []interface{}{priceAttribute})

	return widget
}

// ActionsWidget is used to create an action widget to install the snap.
//
// Returns:
// - Action preview widget for the snap.
func (preview StoreTemplate) ActionsWidget() scopes.PreviewWidget {
	widget := preview.GenericTemplate.ActionsWidget()

	installAction := make(map[string]interface{})
	installAction["id"] = actions.ActionInstall
	installAction["label"] = "Install"

	widget.AddAttributeValue("actions", []interface{}{installAction})

	return widget
}

// UpdatesWidget is used to create a table widget holding snap information.
//
// Returns:
// - Table widget for the snap.
func (preview StoreTemplate) UpdatesWidget() scopes.PreviewWidget {
	widget := preview.GenericTemplate.UpdatesWidget()

	value, ok := widget["values"]
	if ok {
		rows := value.([]interface{})
		if rows != nil {
			sizeRow := []string{"Size", humanize.Bytes(preview.snap.DownloadSize)}
			rows = append(rows, sizeRow)

			widget.AddAttributeValue("values", rows)
		}
	}

	return widget
}
