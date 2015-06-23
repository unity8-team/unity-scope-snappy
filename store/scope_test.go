package store

import (
	//"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical NewScope usage.
func TestNewScope(t *testing.T) {
	scope, err := NewScope(webdm.DefaultApiUrl)
	if err != nil {
		t.Errorf("Unexpected error creating scope: %s", err)
	}

	if scope == nil {
		t.Error("Scope was unexpectedly nil")
	}
}

// Test creating new scope with an invalid API URL
func TestNewScope_invalidUrl(t *testing.T) {
	scope, err := NewScope(":")
	if err == nil {
		t.Errorf("Expected an error creating scope due to invalid URL")
	}

	if scope != nil {
		t.Error("Scope should have been nil")
	}
}

// Test getPackageList for installed packages
func TestGetPackageList_installed(t *testing.T) {
	packageManager := &FakePackageManager{}

	_, err := getPackageList(packageManager, "installed")
	if err != nil {
		t.Error("Unexpected error while getting installed package list")
	}

	if !packageManager.getInstalledPackagesCalled {
		t.Error("Expected GetInstalledPackages() to be called")
	}
}

// Test getPackageList failure getting installed packages
func TestGetPackageList_installed_failure(t *testing.T) {
	packageManager := &FakePackageManager{failToGetInstalledPackages: true}

	packages, err := getPackageList(packageManager, "installed")
	if err == nil {
		t.Error("Expected an error getting installed package list")
	}

	if packages != nil {
		t.Error("Expected no packages to be returned")
	}
}

// Test getPackageList for store packages
func TestGetPackageList_store(t *testing.T) {
	packageManager := &FakePackageManager{}

	_, err := getPackageList(packageManager, "")
	if err != nil {
		t.Error("Unexpected error while getting store package list")
	}

	if !packageManager.getStorePackagesCalled {
		t.Error("Expected GetStorePackages() to be called")
	}
}

// Test getPackageList failure getting store packages
func TestGetPackageList_store_failure(t *testing.T) {
	packageManager := &FakePackageManager{failToGetStorePackages: true}

	packages, err := getPackageList(packageManager, "")
	if err == nil {
		t.Error("Expected an error getting store package list")
	}

	if packages != nil {
		t.Error("Expected no packages to be returned")
	}
}
