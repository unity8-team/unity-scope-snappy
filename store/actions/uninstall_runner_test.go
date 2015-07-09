package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages/fakes"
	"testing"
)

// Test typical Run usage.
func TestUninstallRunner_run(t *testing.T) {
	actionRunner, _ := NewUninstallRunner()

	packageManager := new(fakes.FakeDbusManager)

	response, err := actionRunner.Run(packageManager, "foo")
	if err != nil {
		// Exit here so we don't dereference nil
		t.Fatalf("Unexpected error when attempting to run: %s", err)
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

	if !metadata.UninstallRequested {
		t.Errorf("Expected metadata to indicate that an uninstallation was requested")
	}
}
