/* Copyright (C) 2015 Canonical Ltd.
 *
 * This file is part of unity-scope-snappy.
 *
 * unity-scope-snappy is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option) any
 * later version.
 *
 * unity-scope-snappy is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * unity-scope-snappy. If not, see <http://www.gnu.org/licenses/>.
 */

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
