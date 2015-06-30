package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages/fakes"
	"launchpad.net/unity-scope-snappy/store/progress"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical Run usage.
func TestRefreshInstallingActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewRefreshInstallingRunner()

	packageManager := new(fakes.FakeManager)

	response, err := actionRunner.Run(packageManager, "foo")
	if err != nil {
		// Exit here so we don't dereference nil
		t.Fatalf("Unexpected error when attempting to run: %s", err)
	}

	if response.Status != scopes.ActivationShowPreview {
		t.Errorf(`Response status was "%d", expected "%d"`, response.Status, scopes.ActivationShowPreview)
	}

	// Verify progress hack
	progressHack, ok := response.ScopeData.(progress.Hack)
	if !ok {
		// Exit here so we don't dereference nil
		t.Fatalf("Expected response ScopeData to be a ProgressHack")
	}

	if progressHack.DesiredStatus != webdm.StatusInstalled {
		t.Errorf(`Desired status was "%d", expected "%d"`, progressHack.DesiredStatus, webdm.StatusInstalled)
	}
}
