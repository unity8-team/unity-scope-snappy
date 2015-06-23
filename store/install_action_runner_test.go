package store

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical Run usage.
func TestInstallActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewInstallActionRunner()

	packageManager := new(FakePackageManager)

	response, err := actionRunner.Run(packageManager, "foo")
	if err != nil {
		// Exit here so we don't dereference nil
		t.Fatalf("Unexpected error when attempting to run: %s", err)
	}

	if !packageManager.installCalled {
		t.Error("Expected package manager Install() function to be called")
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

	if progressHack.DesiredStatus != webdm.StatusInstalled {
		t.Errorf(`Desired status was "%d", expected "%d"`, progressHack.DesiredStatus, webdm.StatusInstalled)
	}
}

// Test that a failure to install results in an error
func TestInstallActionRunnerRun_installationFailure(t *testing.T) {
	actionRunner, _ := NewInstallActionRunner()

	packageManager := &FakePackageManager{failToInstall: true}

	response, err := actionRunner.Run(packageManager, "foo")
	if err == nil {
		t.Error("Expected an error due to failure to install")
	}
	if response != nil {
		t.Error("Expected response to be nil")
	}
}
