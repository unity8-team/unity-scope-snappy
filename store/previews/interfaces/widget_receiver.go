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

package interfaces

import "launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"

// WidgetReceiver is an interface to be implemented by any struct that supports
// the type of preview widget interface used by this scope.
type WidgetReceiver interface {
	PushWidgets(widgets ...scopes.PreviewWidget) error
}
