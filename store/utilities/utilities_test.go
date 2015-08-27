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

package utilities

import (
	"launchpad.net/unity-scope-snappy/store/packages/fakes"
	"testing"
)

// Test getPackageList for installed packages
func TestGetPackageList_installed(t *testing.T) {
	packageManager := &fakes.FakeWebdmManager{}

	_, err := GetPackageList(packageManager, "installed", "")
	if err != nil {
		t.Error("Unexpected error while getting installed package list")
	}

	if !packageManager.GetInstalledPackagesCalled {
		t.Error("Expected GetInstalledPackages() to be called")
	}
}

// Test getPackageList failure getting installed packages
func TestGetPackageList_installed_failure(t *testing.T) {
	packageManager := &fakes.FakeWebdmManager{FailGetInstalledPackages: true}

	packages, err := GetPackageList(packageManager, "installed", "")
	if err == nil {
		t.Error("Expected an error getting installed package list")
	}

	if packages != nil {
		t.Error("Expected no packages to be returned")
	}
}

// Test getPackageList for store packages
func TestGetPackageList_store(t *testing.T) {
	packageManager := &fakes.FakeWebdmManager{}

	_, err := GetPackageList(packageManager, "", "")
	if err != nil {
		t.Error("Unexpected error while getting store package list")
	}

	if !packageManager.GetStorePackagesCalled {
		t.Error("Expected GetStorePackages() to be called")
	}
}

// Test getPackageList failure getting store packages
func TestGetPackageList_store_failure(t *testing.T) {
	packageManager := &fakes.FakeWebdmManager{FailGetStorePackages: true}

	packages, err := GetPackageList(packageManager, "", "")
	if err == nil {
		t.Error("Expected an error getting store package list")
	}

	if packages != nil {
		t.Error("Expected no packages to be returned")
	}
}
