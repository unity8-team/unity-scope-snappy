package actions

import (
	"launchpad.net/unity-scope-snappy/store/packages/fakes"
	"testing"
)

// Test typical Run usage.
func TestOpenActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewOpenRunner()

	_, err := actionRunner.Run(&fakes.FakeDbusManager{}, "foo")
	if err == nil {
		t.Error("Expected an error due to opening not being supported")
	}
}
