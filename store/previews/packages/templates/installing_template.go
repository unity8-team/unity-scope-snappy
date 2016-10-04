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
	"github.com/godbus/dbus"
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// InstallingTemplate is a preview template for a package that is currently
// being installed. It's based upon the StoreTemplate.
type InstallingTemplate struct {
	*StoreTemplate
	objectPath dbus.ObjectPath
}

// NewInstallingTemplate creates a new InstallingTemplate.
//
// Parameters:
// snap: Snap to be represented by this template.
// objectPath: DBus object path upon which progress updates will be provided.
//
// Returns:
// - Pointer to new InstallingTemplate (nil if error)
// - Error (nil if none)
func NewInstallingTemplate(snap webdm.Package, objectPath dbus.ObjectPath) (*InstallingTemplate, error) {
	if snap.Uninstalling() {
		return nil, fmt.Errorf("Snap is currently being uninstalled")
	}

	if !objectPath.IsValid() {
		return nil, fmt.Errorf(`Invalid object path: "%s"`, objectPath)
	}

	template := &InstallingTemplate{objectPath: objectPath}

	var err error
	template.StoreTemplate, err = NewStoreTemplate(snap)
	if err != nil {
		return nil, fmt.Errorf("Unable to create store template: %s", err)
	}

	return template, nil
}

// ActionsWidget is used to create a progress widget where the store actions
// were.
//
// Returns:
// - Progress preview widget for the snap.
func (preview InstallingTemplate) ActionsWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("install", "progress")

	source := make(map[string]interface{})
	source["dbus-name"] = "com.canonical.applications.WebdmPackageManager"
	source["dbus-object"] = preview.objectPath

	widget.AddAttributeValue("source", source)

	return widget
}
