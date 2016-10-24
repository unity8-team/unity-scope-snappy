/* Copyright (C) 2016 Canonical Ltd.
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

package daemon

import (
	"testing"
	"time"
)

// Test typical NewSnapdPackageManagerInterface usage.
func TestNewSnapdPackageManagerInterface(t *testing.T) {
	manager, err := NewSnapdPackageManagerInterface(new(FakeDbusServer), "foo", "/foo")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	if manager.progressSignalName != "foo.progress" {
		t.Errorf(`Progress signal name was "%s", expected "foo.progress"`, manager.progressSignalName)
	}

	if manager.finishedSignalName != "foo.finished" {
		t.Errorf(`Finished signal name was "%s", expected "foo.finished"`, manager.finishedSignalName)
	}

	if manager.errorSignalName != "foo.error" {
		t.Errorf(`Error signal name was "%s", expected "foo.error"`, manager.errorSignalName)
	}
}

// Test that NewSnapdmPackageManagerInterface fails with an invalid base object
// path.
func TestNewSnapdPackageManagerInterface_invalidBaseObjectPath(t *testing.T) {
	_, err := NewSnapdPackageManagerInterface(new(FakeDbusServer), "foo", "invalid")
	if err == nil {
		t.Error("Expected an error due to an invalid base object path")
	}
}

// Test typical Install usage.
func TestSnapdInstall(t *testing.T) {
	dbusServer := new(FakeDbusServer)
	dbusServer.InitializeSignals()

	manager, err := NewSnapdPackageManagerInterface(dbusServer, "foo", "/foo")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	// Make the manager poll faster so the tests are more timely
	manager.pollPeriod = time.Millisecond

	// Begin installation of two packages
	_, dbusErr := manager.Install("foo")
	if dbusErr == nil {
		t.Fatalf("Expected error while installing 'foo'")
	}

	_, dbusErr2 := manager.Install("bar")
	if dbusErr2 == nil {
		t.Fatalf("Expected error while installing 'bar'")
	}
}

// Test typical Uninstall usage.
func TestSnapdUninstall(t *testing.T) {
	dbusServer := new(FakeDbusServer)
	dbusServer.InitializeSignals()

	manager, err := NewSnapdPackageManagerInterface(dbusServer, "foo", "/foo")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	// Make the manager poll faster so the tests are more timely
	manager.pollPeriod = time.Millisecond

	// Begin installation of two packages
	_, dbusErr := manager.Uninstall("foo")
	if dbusErr == nil {
		t.Fatalf("Expected error while installing 'foo'")
	}

	_, dbusErr2 := manager.Uninstall("bar")
	if dbusErr2 == nil {
		t.Fatalf("Expected error while installing 'bar'")
	}
}
