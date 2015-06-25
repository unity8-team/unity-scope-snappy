package store

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical Run usage.
func TestUninstallActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewUninstallActionRunner()

	packageManager := new(packages.FakeManager)

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

	// Verify progress hack
	progressHack, ok := response.ScopeData.(ProgressHack)
	if !ok {
		// Exit here so we don't dereference nil
		t.Fatalf("Expected response ScopeData to be a ProgressHack")
	}

	if progressHack.DesiredStatus != webdm.StatusNotInstalled {
		t.Errorf(`Desired status was "%d", expected "%d"`, progressHack.DesiredStatus, webdm.StatusNotInstalled)
	}
}

// Test that a failure to uninstall results in an error
func TestUninstallActionRunnerRun_uninstallationFailure(t *testing.T) {
	actionRunner, _ := NewUninstallActionRunner()

	packageManager := &packages.FakeManager{FailUninstall: true}

	response, err := actionRunner.Run(packageManager, "foo")
	if err == nil {
		t.Error("Expected an error due to failure to uninstall")
	}
	if response != nil {
		t.Error("Unexpected response... expected nil")
	}
}