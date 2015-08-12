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

package mocks

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

// MockBusObject is a mocked implementation of the dbus.BusObject interface, for
// use within tests.
type MockBusObject struct {
	CallCalled bool
	GoCalled bool
	GetPropertyCalled bool
	DestinationCalled bool
	PathCalled bool

	Method string
	Args []interface{}
}

func (mock *MockBusObject) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {
	mock.CallCalled = true
	mock.Method = method
	mock.Args = args
	return &dbus.Call{Body: []interface{}{dbus.ObjectPath("/foo/1")}}
}

func (mock *MockBusObject) Go(method string, flags dbus.Flags, ch chan *dbus.Call, args ...interface{}) *dbus.Call {
	mock.GoCalled = true
	mock.Method = method
	mock.Args = args
	return &dbus.Call{Body: []interface{}{dbus.ObjectPath("/foo/1")}}
}

func (mock *MockBusObject) GetProperty(p string) (dbus.Variant, error) {
	mock.GetPropertyCalled = true
	return dbus.MakeVariant("foo"), nil
}

func (mock *MockBusObject) Destination() string {
	mock.DestinationCalled = true
	return "foo"
}

func (mock *MockBusObject) Path() dbus.ObjectPath {
	mock.PathCalled = true
	return dbus.ObjectPath("/foo/1")
}
