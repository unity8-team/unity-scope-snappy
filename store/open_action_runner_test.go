package store

import (
	"launchpad.net/unity-scope-snappy/store/packages"
	"testing"
)

// Test typical Run usage.
func TestOpenActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewOpenActionRunner()

	packageManager := new(packages.FakeManager)

	_, err := actionRunner.Run(packageManager, "foo")
	if err == nil {
		t.Error("Expected an error due to opening not being supported")
	}
}
