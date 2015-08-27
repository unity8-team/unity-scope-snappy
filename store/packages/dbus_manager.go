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
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

// DbusManager is an interface to be implemented by any struct that supports
// asynchronous package installation/uninstallation via dbus in the manner used
// by this scope.
type DbusManager interface {
	Connect() error
	Install(packageId string) (dbus.ObjectPath, error)
	Uninstall(packageId string) (dbus.ObjectPath, error)
}
