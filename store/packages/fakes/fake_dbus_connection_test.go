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
	"launchpad.net/unity-scope-snappy/store/packages/mocks"
	"reflect"
	"testing"
)

// Test that Names returns a single name
func TestFakeDbusConnection_names(t *testing.T) {
	connection := FakeDbusConnection{}

	names := connection.Names()
	if len(names) != 1 {
		t.Fatalf("Got %d names, expected 1", len(names))
	}
}

// Test that Object simply returns the given DbusObject.
func TestFakeDbusConnection_object(t *testing.T) {
	mock := &mocks.MockBusObject{}
	connection := FakeDbusConnection{mock}

	if !reflect.DeepEqual(connection.Object("foo", "bar"), mock) {
		t.Error("Expected the fake connection to return the given mock")
	}
}
