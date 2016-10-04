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
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// GenericTemplate is a Template implementation that doesn't contain any
// conditionals depending on package information. It's meant to be embedded in
// other structs and further specialized.
type GenericTemplate struct {
	snap webdm.Package
}

// NewGenericTemplate creates a new GenericTemplate.
func NewGenericTemplate(snap webdm.Package) *GenericTemplate {
	return &GenericTemplate{snap: snap}
}

// HeaderWidget is used to create a header widget for the snap.
//
// Returns:
// - Header preview widget for the snap.
func (preview GenericTemplate) HeaderWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("header", "header")

	widget.AddAttributeMapping("title", "title")
	widget.AddAttributeMapping("subtitle", "subtitle")
	widget.AddAttributeMapping("mascot", "art")

	return widget
}

// ActionsWidget is used to create an action widget for the snap. The widget
// contains no actions.
//
// Returns:
// - Empty action preview widget for the snap.
func (preview GenericTemplate) ActionsWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("actions", "actions")

	return widget
}

// InfoWidget is used to create a text widget holding the snap description.
//
// Returns:
// - Text preview widget for the snap.
func (preview GenericTemplate) InfoWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("summary", "text")

	widget.AddAttributeValue("title", "Info")
	widget.AddAttributeValue("text", preview.snap.Description)

	return widget
}

// UpdatesWidget is used to create a table widget holding snap version
// information.
//
// Returns:
// - Table widget for the snap.
func (preview GenericTemplate) UpdatesWidget() scopes.PreviewWidget {
	widget := scopes.NewPreviewWidget("updates_table", "table")
	widget.AddAttributeValue("title", "Updates")

	versionRow := []string{"Version number", preview.snap.Version}

	widget.AddAttributeValue("values", []interface{}{versionRow})

	return widget
}
