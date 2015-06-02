package main

import (
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical Run usage.
func TestOkActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewOkActionRunner()

	packageManager := new(FakePackageManager)

	response, err := actionRunner.Run(packageManager, "foo")
	if err != nil {
		// Exit here so we don't dereference nil
		t.Fatalf("Unexpected error when attempting to run: %s", err)
	}

	if response.Status != scopes.ActivationShowPreview {
		t.Errorf(`Response status was "%d", expected "%d"`, response.Status, scopes.ActivationShowPreview)
	}

	// Verify no progress hack
	progressHack, ok := response.ScopeData.(ProgressHack)
	if ok {
		if progressHack.DesiredStatus != webdm.StatusUndefined {
			t.Errorf("No progress hack should be included")
		}
	}
}
