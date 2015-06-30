package actions

import (
	"launchpad.net/unity-scope-snappy/store/packages/fakes"
	"testing"
)

// Test typical Run usage.
func TestOpenActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewOpenRunner()

	packageManager := new(fakes.FakeManager)

	_, err := actionRunner.Run(packageManager, "foo")
	if err == nil {
		t.Error("Expected an error due to opening not being supported")
	}
}
