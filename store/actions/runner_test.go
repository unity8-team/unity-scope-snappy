package actions

import (
	"reflect"
	"testing"
)

// Data for TestNewRunner
var newRunnerTests = []struct {
	actionId ActionId
	expected interface{}
}{
	{ActionInstall, &InstallRunner{}},
	{ActionUninstall, &UninstallRunner{}},
	{ActionUninstallConfirm, &ConfirmUninstallRunner{}},
	{ActionUninstallCancel, &CancelUninstallRunner{}},
	{ActionOpen, &OpenRunner{}},
	{ActionFinished, &FinishedRunner{}},
	{ActionFailed, &FailedRunner{}},
}

// Test typical NewRunner usage.
func TestNewRunner(t *testing.T) {
	for i, test := range newRunnerTests {
		runner, err := NewRunner(test.actionId)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error: %s", i, err)
		} else {
			runnerType := reflect.TypeOf(runner)
			expectedType := reflect.TypeOf(test.expected)
			if runnerType != expectedType {
				t.Errorf(`Test case %d: Action runner type was "%s", expected "%s"`,
					i, runnerType, expectedType)
			}
		}
	}
}

// Test that an invalid action ID results in an error
func TestNewRunner_invalidAction(t *testing.T) {
	runner, err := NewRunner(ActionId(0))
	if err == nil {
		t.Error("Expected an error due to invalid action ID")
	}

	if runner != nil {
		t.Error("Expected action runner to be nil due to error")
	}
}
