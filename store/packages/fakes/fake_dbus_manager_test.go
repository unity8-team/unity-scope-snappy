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
	"testing"
)

// Test typical Connect usage.
func TestFakeDbusManager_connect(t *testing.T) {
	manager := &FakeDbusManager{}

	err := manager.Connect()
	if err != nil {
		t.Fatalf("Unexpected error while connecting: %s", err)
	}

	if !manager.ConnectCalled {
		t.Error("Expected ConnectCalled to have been set")
	}
}

// Test that requesting an error in Connect actually results in an error.
func TestFakeDbusManager_connect_failureRequested(t *testing.T) {
	manager := &FakeDbusManager{FailConnect: true}

	err := manager.Connect()
	if err == nil {
		t.Error("Expected an error due to failure request")
	}

	if !manager.ConnectCalled {
		t.Error("Expected ConnectCalled to have been set")
	}
}

// Test typical Install usage.
func TestFakeDbusManager_Install(t *testing.T) {
	manager := &FakeDbusManager{}

	objectPath, err := manager.Install("foo")
	if err != nil {
		t.Fatalf("Unexpected error while installing: %s", err)
	}

	if !objectPath.IsValid() {
		t.Errorf("Object path was unexpectedly invalid: %s", objectPath)
	}

	if !manager.InstallCalled {
		t.Error("Expected InstallCalled to have been set")
	}
}

// Test that requesting an error in Install actually results in an error.
func TestFakeDbusManager_Install_failureRequest(t *testing.T) {
	manager := &FakeDbusManager{FailInstall: true}

	_, err := manager.Install("foo")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}

	if !manager.InstallCalled {
		t.Error("Expected InstallCalled to have been set")
	}
}

// Test typical Uninstall usage.
func TestFakeDbusManager_Uninstall(t *testing.T) {
	manager := &FakeDbusManager{}

	objectPath, err := manager.Uninstall("foo")
	if err != nil {
		t.Fatalf("Unexpected error while uninstalling: %s", err)
	}

	if !objectPath.IsValid() {
		t.Errorf("Object path was unexpectedly invalid: %s", objectPath)
	}

	if !manager.UninstallCalled {
		t.Error("Expected UninstallCalled to have been set")
	}
}

// Test that requesting an error in Uninstall actually results in an error.
func TestFakeDbusManager_Uninstall_failureRequest(t *testing.T) {
	manager := &FakeDbusManager{FailUninstall: true}

	_, err := manager.Uninstall("foo")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}

	if !manager.UninstallCalled {
		t.Error("Expected UninstallCalled to have been set")
	}
}
