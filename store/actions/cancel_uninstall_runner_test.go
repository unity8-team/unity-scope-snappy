package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages/fakes"
	"testing"
)

// Test typical Run usage.
func TestCancelUninstallRunner_run(t *testing.T) {
	runner, _ := NewCancelUninstallRunner()

	response, err := runner.Run(&fakes.FakeDbusManager{}, "foo")
	if err != nil {
		// Exit here so we don't dereference nil
		t.Fatalf("Unexpected error when attempting to run: %s", err)
	}

	if response.Status != scopes.ActivationShowPreview {
		t.Errorf(`Response status was "%d", expected "%d"`, response.Status, scopes.ActivationShowPreview)
	}

	// Verify lack of operation metadata
	_, ok := response.ScopeData.(operation.Metadata)
	if ok {
		t.Error("Response ScopeData should not include operation metadata")
	}
}
