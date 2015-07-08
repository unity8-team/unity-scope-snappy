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
