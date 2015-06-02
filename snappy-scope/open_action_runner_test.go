package main

import (
	"testing"
)

// Test typical Run usage.
func TestOpenActionRunnerRun(t *testing.T) {
	actionRunner, _ := NewOpenActionRunner()

	packageManager := new(FakePackageManager)

	_, err := actionRunner.Run(packageManager, "foo")
	if err == nil {
		t.Error("Expected an error due to opening not being supported")
	}
}
