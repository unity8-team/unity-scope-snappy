package scope

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical Run usage.
func TestRefreshUninstallingActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewRefreshUninstallingActionRunner()

	packageManager := new(FakePackageManager)

	response, err := actionRunner.Run(packageManager, "foo")
	if err != nil {
		// Exit here so we don't dereference nil
		t.Fatalf("Unexpected error when attempting to run: %s", err)
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