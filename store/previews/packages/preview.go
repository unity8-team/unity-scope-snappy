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

package packages

import (
	"github.com/snapcore/snapd/client"
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
	"launchpad.net/unity-scope-snappy/store/previews/packages/templates"
)

// Preview is a PreviewGenerator representing a given package.
type Preview struct {
	template templates.Template
}

// NewPreview creates a new Preview for representing a given package.
//
// Parameters:
// snap: Package to be represented by the preview.
func NewPreview(snap client.Snap, result *scopes.Result, metadata operation.Metadata) (*Preview, error) {
	preview := new(Preview)
	var err error

	installed := false
	if (snap.Status == client.StatusInstalled ||
		snap.Status == client.StatusActive) {
		installed = true
	}
	if metadata.InstallRequested  && !installed {
		preview.template, err = templates.NewInstallingTemplate(snap, result, metadata.ObjectPath)
	} else if metadata.UninstallConfirmed && installed {
		preview.template, err = templates.NewUninstallingTemplate(snap, metadata.ObjectPath)
	} else {
		if installed {
			preview.template, err = templates.NewInstalledTemplate(snap)
		} else {
			preview.template, err = templates.NewStoreTemplate(snap, result)
		}
	}

	return preview, err
}

// Generate pushes the template's preview widgets onto a WidgetReceiver.
//
// Parameters:
// receiver: Implementation of the WidgetReceiver interface.
//
// Returns:
// - Error (nil if none)
func (preview Preview) Generate(receiver interfaces.WidgetReceiver) error {
	receiver.PushWidgets(preview.template.HeaderWidget())
	receiver.PushWidgets(preview.template.ActionsWidget())
	receiver.PushWidgets(preview.template.InfoWidget())
	receiver.PushWidgets(preview.template.UpdatesWidget())

	return nil
}
