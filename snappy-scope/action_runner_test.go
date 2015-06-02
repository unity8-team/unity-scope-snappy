package main

import (
	"reflect"
	"testing"
)

// Data for TestNewActionRunner
var newActionRunnerTests = []struct {
	actionId     ActionId
	expectedType string
}{
	{ActionInstall, "*main.InstallActionRunner"},
	{ActionUninstall, "*main.UninstallActionRunner"},
	{ActionOpen, "*main.OpenActionRunner"},

	// Temporary actions for manual refresh
	{ActionRefreshInstalling, "*main.RefreshInstallingActionRunner"},
	{ActionRefreshUninstalling, "*main.RefreshUninstallingActionRunner"},
	{ActionOk, "*main.OkActionRunner"},
}

// Test typical NewPreview usage.
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
