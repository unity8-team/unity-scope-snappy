package actions

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages/fakes"
	"testing"
)

// Test typical Run usage.
func TestConfirmUninstallRunner_run(t *testing.T) {
	actionRunner, _ := NewConfirmUninstallRunner()

	packageManager := new(fakes.FakeDbusManager)

	response, err := actionRunner.Run(packageManager, "foo")
	if err != nil {
		// Exit here so we don't dereference nil
		t.Fatalf("Unexpected error when attempting to run: %s", err)
	}

	if !packageManager.UninstallCalled {
		t.Error("Expected package manager Uninstall() function to be called")
	}

	if response.Status != scopes.ActivationShowPreview {
		t.Errorf(`Response status was "%d", expected "%d"`, response.Status, scopes.ActivationShowPreview)
	}

	// Verify operation metadata
	metadata, ok := response.ScopeData.(operation.Metadata)
	if !ok {
		// Exit here so we don't dereference nil
		t.Fatalf("Expected response ScopeData to include operation metadata")
	}

	if !metadata.UninstallConfirmed {
		t.Errorf("Expected metadata to indicate that an uninstallation was confirmed")
	}

	if metadata.ObjectPath != dbus.ObjectPath("/foo/1") {
		t.Errorf(`Metadata object path was "%s", expected "/foo/1"`, metadata.ObjectPath)
	}
}

// Test that a failure to uninstall results in an error
func TestConfirmUninstallRunner_run_uninstallationFailure(t *testing.T) {
	actionRunner, _ := NewConfirmUninstallRunner()

	packageManager := &fakes.FakeDbusManager{FailUninstall: true}

	response, err := actionRunner.Run(packageManager, "foo")
	if err == nil {
		t.Error("Expected an error due to failure to uninstall")
	}
	if response != nil {
		t.Error("Unexpected response... expected nil")
	}
}
