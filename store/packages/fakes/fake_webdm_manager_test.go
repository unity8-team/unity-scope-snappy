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

// Test typical GetInstalledPackages usage.
func TestFakeWebdmManager_GetInstalledPackages(t *testing.T) {
	manager := &FakeWebdmManager{}

	packages, err := manager.GetInstalledPackages("")
	if err != nil {
		t.Fatalf("Unexpected error while getting installed packages: %s", err)
	}

	if len(packages) < 1 {
		t.Errorf("Got %d packages, expected at least 1", len(packages))
	}

	if !manager.GetInstalledPackagesCalled {
		t.Error("Expected GetInstalledPackagesCalled to have been set")
	}
}

// Test that requesting an error in GetInstalledPackages actually results in an
// error.
func TestFakeWebdmManager_GetInstalledPackages_failureRequest(t *testing.T) {
	manager := &FakeWebdmManager{FailGetInstalledPackages: true}

	_, err := manager.GetInstalledPackages("")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}

	if !manager.GetInstalledPackagesCalled {
		t.Error("Expected GetInstalledPackagesCalled to have been set")
	}
}

// Test typical GetStorePackages usage.
func TestFakeWebdmManager_GetStorePackages(t *testing.T) {
	manager := &FakeWebdmManager{}

	packages, err := manager.GetStorePackages("")
	if err != nil {
		t.Fatalf("Unexpected error while getting store packages: %s", err)
	}

	if len(packages) < 1 {
		t.Errorf("Got %d packages, expected at least 1", len(packages))
	}

	if !manager.GetStorePackagesCalled {
		t.Error("Expected GetStorePackagesCalled to have been set")
	}
}

// Test that requesting an error in GetStorePackages actually results in an
// error.
func TestFakeWebdmManager_GetStorePackages_failureRequest(t *testing.T) {
	manager := &FakeWebdmManager{FailGetStorePackages: true}

	_, err := manager.GetStorePackages("")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}

	if !manager.GetStorePackagesCalled {
		t.Error("Expected GetStorePackagesCalled to have been set")
	}
}

// Test typical Query usage.
func TestFakeWebdmManager_Query(t *testing.T) {
	manager := &FakeWebdmManager{}

	snap, err := manager.Query("foo")
	if err != nil {
		t.Fatalf("Unexpected error while querying: %s", err)
	}

	if snap == nil {
		t.Error("Snap was unexpectedly nil")
	}

	if !manager.QueryCalled {
		t.Error("Expected QueryCalled to have been set")
	}
}

// Test that requesting an error in Query actually results in an error.
func TestFakeWebdmManager_Query_failureRequest(t *testing.T) {
	manager := &FakeWebdmManager{FailQuery: true}

	_, err := manager.Query("foo")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}

	if !manager.QueryCalled {
		t.Error("Expected QueryCalled to have been set")
	}
}

// Test typical Install usage.
func TestFakeWebdmManager_Install(t *testing.T) {
	manager := &FakeWebdmManager{}

	err := manager.Install("foo")
	if err != nil {
		t.Fatalf("Unexpected error while installing: %s", err)
	}

	if !manager.InstallCalled {
		t.Error("Expected InstallCalled to have been set")
	}
}

// Test that requesting an error in Install actually results in an error.
func TestFakeWebdmManager_Install_failureRequest(t *testing.T) {
	manager := &FakeWebdmManager{FailInstall: true}

	err := manager.Install("foo")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}

	if !manager.InstallCalled {
		t.Error("Expected InstallCalled to have been set")
	}
}

// Test typical Uninstall usage.
func TestFakeWebdmManager_Uninstall(t *testing.T) {
	manager := &FakeWebdmManager{}

	err := manager.Uninstall("foo")
	if err != nil {
		t.Fatalf("Unexpected error while uninstalling: %s", err)
	}

	if !manager.UninstallCalled {
		t.Error("Expected UninstallCalled to have been set")
	}
}

// Test that requesting an error in Uninstall actually results in an error.
func TestFakeWebdmManager_Uninstall_failureRequest(t *testing.T) {
	manager := &FakeWebdmManager{FailUninstall: true}

	err := manager.Uninstall("foo")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}

	if !manager.UninstallCalled {
		t.Error("Expected UninstallCalled to have been set")
	}
}
