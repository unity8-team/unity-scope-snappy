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

package previews

import (
	"fmt"

	"github.com/snapcore/snapd/client"
	"launchpad.net/unity-scope-snappy/store/actions"
	"launchpad.net/unity-scope-snappy/store/previews/fakes"
	"testing"
)

// Test typical NewPackagePreview usage.
func TestNewConfirmUninstallPreview(t *testing.T) {
	snap := client.Snap{Name: "package1"}

	preview := NewConfirmUninstallPreview(snap)
	if preview == nil {
		t.Fatal("Preview was unexpectedly nil")
	}

	if preview.snap.Name != snap.Name {
		t.Error(`Preview snap name was "%s", expected "%s"`, preview.snap.Name,
			snap.Name)
	}
}

// Test typical Generate usage, and verify that it conforms to store design.
func TestConfirmUninstallPreview_generate(t *testing.T) {
	snap := client.Snap{Name: "package1"}
	preview := NewConfirmUninstallPreview(snap)

	receiver := new(fakes.FakeWidgetReceiver)

	err := preview.Generate(receiver)
	if err != nil {
		t.Errorf("Unexpected error while generating preview: %s", err)
	}

	if len(receiver.Widgets) != 2 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Got %d widgets, expected 2", len(receiver.Widgets))
	}

	// Verify text
	widget := receiver.Widgets[0]
	if widget.WidgetType() != "text" {
		t.Error("Expected text to be first widget")
	}

	value, ok := widget["text"]
	if !ok {
		t.Error(`Expected text widget to contain "text"`)
	}

	expectedText := fmt.Sprintf("Are you sure you want to uninstall %s?", snap.Name)
	if value != expectedText {
		t.Errorf(`Text was "%s", expected "%s"`, value, expectedText)
	}

	// Verify actions
	widget = receiver.Widgets[1]
	if widget.WidgetType() != "actions" {
		t.Fatal("Expected actions to be second widget")
	}

	value, ok = widget["actions"]
	if !ok {
		t.Fatal(`Expected actions widget to include "actions"`)
	}

	actionsInterfaces := value.([]interface{})

	if len(actionsInterfaces) != 2 {
		t.Fatalf("Actions widget had %d actions, expected 2", len(actionsInterfaces))
	}

	// Verify the uninstall action
	action := actionsInterfaces[0].(map[string]interface{})
	value, ok = action["id"]
	if !ok {
		t.Errorf("Expected uninstall action to have an id")
	}
	if value != actions.ActionUninstallConfirm {
		t.Errorf(`Open action's ID was "%s", expected "%s"`, value, actions.ActionUninstallConfirm)
	}

	value, ok = action["label"]
	if !ok {
		t.Errorf("Expected open action to have a label")
	}
	if value != "Uninstall" {
		t.Errorf(`Open action's label was "%s", expected "Uninstall"`, value)
	}

	// Verify the cancel action
	action = actionsInterfaces[1].(map[string]interface{})
	value, ok = action["id"]
	if !ok {
		t.Errorf("Expected cancel action to have an id")
	}
	if value != actions.ActionUninstallCancel {
		t.Errorf(`Cancel action's ID was "%s", expected "%s"`, value, actions.ActionUninstallCancel)
	}

	value, ok = action["label"]
	if !ok {
		t.Errorf("Expected cancel action to have a label")
	}
	if value != "Cancel" {
		t.Errorf(`Cancel action's label was "%s", expected "Cancel"`, value)
	}
}
