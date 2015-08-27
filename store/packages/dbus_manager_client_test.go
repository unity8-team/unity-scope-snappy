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
	"testing"
)

// Test typical NewDbusManagerClient usage.
func TestNewDbusManagerClient(t *testing.T) {
	client := NewDbusManagerClient()
	if client == nil {
		t.Fatal("Client was unexpectedly nil")
	}

	if client.dbusObject == "" {
		t.Error("Client didn't configure a dbus object")
	}

	if client.dbusObjectInterface == "" {
		t.Error("Client didn't configure a dbus object interface")
	}
}

// Test typical Connect usage.
func TestDbusManagerClient_connect(t *testing.T) {
	client := NewDbusManagerClient()
	err := client.Connect()
	if err != nil {
		t.Errorf("Unexpected error while connecting: %s", err)
	}

	names := client.connection.Names()
	if len(names) < 1 {
		t.Errorf("Got %d names, expected at least 1", len(names))
	}
}

// Test that trying to install before connecting results in an error.
func TestDbusManagerClient_install_beforeConnect(t *testing.T) {
	client := NewDbusManagerClient()
	_, err := client.Install("foo")
	if err == nil {
		t.Error("Expected an error due to install before connect")
	}
}

// Test that trying to uninstall before connecting results in an error.
func TestDbusManagerClient_uninstall_beforeConnect(t *testing.T) {
	client := NewDbusManagerClient()
	_, err := client.Uninstall("foo")
	if err == nil {
		t.Error("Expected an error due to uninstall before connect")
	}
}
