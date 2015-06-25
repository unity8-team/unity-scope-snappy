package packages

import (
	"testing"
)

// Test typical GetInstalledPackages usage.
func TestGetInstalledPackages(t *testing.T) {
	manager := &FakeManager{}

	packages, err := manager.GetInstalledPackages()
	if err != nil {
		t.Fatalf("Unexpected error while getting installed packages: %s", err)
	}

	if len(packages) < 1 {
		t.Errorf("Got %d packages, expected at least 1", len(packages))
	}
}

// Test that requesting an error in GetInstalledPackages actually results in an
// error.
func TestGetInstalledPackages_failureRequest(t *testing.T) {
	manager := &FakeManager{FailGetInstalledPackages: true}

	_, err := manager.GetInstalledPackages()
	if err == nil {
		t.Error("Expected an error due to failure request")
	}
}

// Test typical GetStorePackages usage.
func TestGetStorePackages(t *testing.T) {
	manager := &FakeManager{}

	packages, err := manager.GetStorePackages()
	if err != nil {
		t.Fatalf("Unexpected error while getting store packages: %s", err)
	}

	if len(packages) < 1 {
		t.Errorf("Got %d packages, expected at least 1", len(packages))
	}
}

// Test that requesting an error in GetStorePackages actually results in an
// error.
func TestGetStorePackages_failureRequest(t *testing.T) {
	manager := &FakeManager{FailGetStorePackages: true}

	_, err := manager.GetStorePackages()
	if err == nil {
		t.Error("Expected an error due to failure request")
	}
}

// Test typical Query usage.
func TestQuery(t *testing.T) {
	manager := &FakeManager{}

	snap, err := manager.Query("foo")
	if err != nil {
		t.Fatalf("Unexpected error while querying: %s", err)
	}

	if snap == nil {
		t.Error("Snap was unexpectedly nil")
	}
}

// Test that requesting an error in Query actually results in an error.
func TestQuery_failureRequest(t *testing.T) {
	manager := &FakeManager{FailQuery: true}

	_, err := manager.Query("foo")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}
}

// Test typical Install usage.
func TestInstall(t *testing.T) {
	manager := &FakeManager{}

	err := manager.Install("foo")
	if err != nil {
		t.Fatalf("Unexpected error while installing: %s", err)
	}
}

// Test that requesting an error in Install actually results in an error.
func TestInstall_failureRequest(t *testing.T) {
	manager := &FakeManager{FailInstall: true}

	err := manager.Install("foo")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}
}

// Test typical Uninstall usage.
func TestUninstall(t *testing.T) {
	manager := &FakeManager{}

	err := manager.Uninstall("foo")
	if err != nil {
		t.Fatalf("Unexpected error while uninstalling: %s", err)
	}
}

// Test that requesting an error in Uninstall actually results in an error.
func TestUninstall_failureRequest(t *testing.T) {
	manager := &FakeManager{FailUninstall: true}

	err := manager.Uninstall("foo")
	if err == nil {
		t.Error("Expected an error due to failure request")
	}
}
