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
	"fmt"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
	"launchpad.net/unity-scope-snappy/store/previews/packages/templates"
	"launchpad.net/unity-scope-snappy/webdm"
)

// Preview is a PreviewGenerator representing a given package.
type Preview struct {
	template templates.Template
}

// NewPreview creates a new Preview for representing a given package.
//
// Parameters:
// snap: Package to be represented by the preview.
func NewPreview(snap webdm.Package, metadata operation.Metadata) (*Preview, error) {
	preview := new(Preview)
	var err error

	if metadata.InstallRequested && !snap.Installed() {
		if snap.Uninstalling() {
			return nil, fmt.Errorf("Install requested, but package is uninstalling")
		}

		preview.template, err = templates.NewInstallingTemplate(snap, metadata.ObjectPath)
	} else if metadata.UninstallConfirmed && !snap.NotInstalled() {
		if snap.Installing() {
			return nil, fmt.Errorf("Uninstall requested, but package is installing")
		}

		preview.template, err = templates.NewUninstallingTemplate(snap, metadata.ObjectPath)
	} else {
		if snap.Installed() {
			preview.template, err = templates.NewInstalledTemplate(snap)
		} else {
			preview.template, err = templates.NewStoreTemplate(snap)
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
