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

package fakes

import (
	"github.com/godbus/dbus"
)

// DbusConnection is an interface of the DBus connection used by the DBus
// package manager client.
type FakeDbusConnection struct {
	DbusObject dbus.BusObject
}

func (fake FakeDbusConnection) Names() []string {
	return []string{":1.42"}
}

func (fake FakeDbusConnection) Object(dest string, path dbus.ObjectPath) dbus.BusObject {
	return fake.DbusObject
}
