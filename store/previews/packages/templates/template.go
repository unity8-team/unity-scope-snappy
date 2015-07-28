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
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
)

// Template is an interface to be implemented by structs which
// are representing a given package for a Unity scope preview.
type Template interface {
	// HeaderWidget generates a widget for the preview header section.
	HeaderWidget() scopes.PreviewWidget

	// ActionsWidget generates a widget for the preview actions section.
	ActionsWidget() scopes.PreviewWidget

	// InfoWidget generates a widget for the preview info section.
	InfoWidget() scopes.PreviewWidget

	// UpdatesWidget generates a widget for the preview updates section.
	UpdatesWidget() scopes.PreviewWidget
}
