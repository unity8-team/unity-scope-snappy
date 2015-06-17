package scope

import (
	"reflect"
	"testing"
)

// Data for TestNewActionRunner
var newActionRunnerTests = []struct {
	actionId     ActionId
	expectedType string
}{
	{ActionInstall, "*scope.InstallActionRunner"},
	{ActionUninstall, "*scope.UninstallActionRunner"},
	{ActionOpen, "*scope.OpenActionRunner"},

	// Temporary actions for manual refresh
	{ActionRefreshInstalling, "*scope.RefreshInstallingActionRunner"},
	{ActionRefreshUninstalling, "*scope.RefreshUninstallingActionRunner"},
	{ActionOk, "*scope.OkActionRunner"},
}

// Test typical NewActionRunner usage.
func TestNewActionRunner(t *testing.T) {
	for i, test := range newActionRunnerTests {
		actionRunner, err := NewActionRunner(test.actionId)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error: %s", i, err)
		} else {
			actionRunnerType := reflect.TypeOf(actionRunner)
			if actionRunnerType.String() != test.expectedType {
				t.Errorf(`Test case %d: Action runner type was "%s", expected "%s"`,
					i, actionRunnerType, test.expectedType)
			}
		}
	}
}

// Test that an invalid action ID results in an error
func TestNewActionRunning_invalidAction(t *testing.T) {
	actionRunner, err := NewActionRunner(ActionId(0))
	if err == nil {
		t.Error("Expected an error due to invalid action ID")
	}

	if actionRunner != nil {
		t.Error("Expected action runner to be nil due to error")
	}
}
