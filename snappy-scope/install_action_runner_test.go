package main

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
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
}
