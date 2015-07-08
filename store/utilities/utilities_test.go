package utilities

import (
	"launchpad.net/unity-scope-snappy/store/packages/fakes"
	"testing"
)

// Test getPackageList for installed packages
func TestGetPackageList_installed(t *testing.T) {
	packageManager := &fakes.FakeWebdmManager{}

	_, err := GetPackageList(packageManager, "installed")
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

	packages, err := GetPackageList(packageManager, "installed")
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

	_, err := GetPackageList(packageManager, "")
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

	packages, err := GetPackageList(packageManager, "")
	if err == nil {
		t.Error("Expected an error getting store package list")
	}

	if packages != nil {
		t.Error("Expected no packages to be returned")
	}
}
