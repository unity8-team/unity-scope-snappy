package scope

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"testing"
)

// Test typical Run usage.
func TestUninstallActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewUninstallActionRunner()

	packageManager := new(FakePackageManager)

	response, err := actionRunner.Run(packageManager, "foo")
	if err != nil {
		// Exit here so we don't dereference nil
		t.Fatalf("Unexpected error when attempting to run: %s", err)
	}

	if !packageManager.uninstallCalled {
		t.Error("Expected package manager Uninstall() function to be called")
	}

	if response.Status != scopes.ActivationShowPreview {
		t.Errorf(`Response status was "%d", expected "%d"`, response.Status, scopes.ActivationShowPreview)
	}
}

// Test that a failure to uninstall results in an error
func TestUninstallActionRunnerRun_uninstallationFailure(t *testing.T) {
	actionRunner, _ := NewUninstallActionRunner()

	packageManager := &FakePackageManager{failToUninstall: true}

	response, err := actionRunner.Run(packageManager, "foo")
	if err == nil {
		t.Error("Expected an error due to failure to uninstall")
	}
	if response != nil {
		t.Error("Unexpected response... expected nil")
	}
}
