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

package templates

import (
	"github.com/snapcore/snapd/client"
	"launchpad.net/unity-scope-snappy/store/actions"
	"testing"
)

// Data for InstalledTemplate tests
var installedTemplateTests = []struct {
	snap client.Snap
}{
	{client.Snap{ID: "package1", Status: client.StatusInstalled, Version: "0.1", InstalledSize: 123456}},
	{client.Snap{ID: "package1", Status: client.StatusActive, Version: "0.1", InstalledSize: 123456}},
}

// Test typical NewInstalledTemplate usage.
func TestNewInstalledTemplate(t *testing.T) {
	for i, test := range installedTemplateTests {
		template, err := NewInstalledTemplate(test.snap)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		if template.snap.ID != test.snap.ID {
			t.Errorf(`Test case %d: Template snap's ID is "%s", expected "%s"`, i, template.snap.ID, test.snap.ID)
		}
	}
}

// Test that the header widget conforms to the store design.
func TestInstalledTemplate_headerWidget(t *testing.T) {
	for i, test := range installedTemplateTests {
		template, err := NewInstalledTemplate(test.snap)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		widget := template.HeaderWidget()

		// Check generic attributes
		value, ok := widget["attributes"]
		if !ok {
			t.Errorf("Test case %d: Expected header attributes to include generic attributes", i)
			continue
		}

		attributes := value.([]interface{})
		if len(attributes) != 1 {
			t.Errorf("Test case %d: Got %d generic attributes for header, expected 1", i, len(attributes))
			continue
		}

		attribute := attributes[0].(map[string]interface{})
		value, ok = attribute["value"]
		if !ok {
			t.Errorf(`Test case %d: Expected generic header attribute to have "value" key`, i)
		}
		if value != "✔ INSTALLED" {
			t.Errorf(`Test case %d: Generic header attribute "value" was "%s", expected "✔ INSTALLED"`, i, value)
		}
	}
}

// Test that the actions widget conforms to the store design.
func TestInstalledTemplate_actionsWidget(t *testing.T) {
	for i, test := range installedTemplateTests {
		template, err := NewInstalledTemplate(test.snap)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		widget := template.ActionsWidget()

		value, ok := widget["actions"]
		if !ok {
			t.Errorf("Test case %d: Expected actions widget to include actions", i)
			continue
		}

		actionsInterfaces := value.([]interface{})

		// Can only test for nil result, so no Open button
		if len(actionsInterfaces) != 1 {
			t.Errorf("Test case %d: Actions widget has %d actions, expected 2", i, len(actionsInterfaces))
			continue
		}

		// Verify the open action
		// FIXME: Open action cannot currently be tested
		action := actionsInterfaces[0].(map[string]interface{})
/*
		value, ok = action["id"]
		if !ok {
			t.Errorf("Test case %d: Expected open action to have an id", i)
		}
		if value != actions.ActionOpen {
			t.Errorf(`Test case %d: Open action's ID was "%s", expected "%s"`, i, value, actions.ActionOpen)
		}

		value, ok = action["label"]
		if !ok {
			t.Errorf("Test case %d: Expected open action to have a label", i)
		}
		if value != "Open" {
			t.Errorf(`Test case %d: Open action's label was "%s", expected "Open"`, i, value)
		}
*/

		// Verify the uninstall action
		// action = actionsInterfaces[1].(map[string]interface{})
		value, ok = action["id"]
		if !ok {
			t.Errorf("Test case %d: Expected uninstall action to have an id", i)
		}
		if value != actions.ActionUninstall {
			t.Errorf(`Test case %d: Uninstall action's ID was "%s", expected "%s"`, i, value, actions.ActionUninstall)
		}

		value, ok = action["label"]
		if !ok {
			t.Errorf("Test case %d: Expected uninstall action to have a label", i)
		}
		if value != "Uninstall" {
			t.Errorf(`Test case %d: Uninstall action's label was "%s", expected "Uninstall"`, i, value)
		}
	}
}

// Test that the updates widget conforms to the store design.
func TestTestInstalledTemplate_updatesWidget(t *testing.T) {
	for i, test := range installedTemplateTests {
		template, err := NewInstalledTemplate(test.snap)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		widget := template.UpdatesWidget()

		// Verify title
		value, ok := widget["title"]
		if !ok {
			t.Errorf("Test case %d: Expected updates table to include a title", i)
		}
		if value != "Updates" {
			t.Errorf(`Test case %d: Updates table's title was "%s", expected "Updates"`, i, value)
		}

		// Verify table rows
		value, ok = widget["values"]
		if !ok {
			t.Errorf("Test case %d: Expected updates table to include values", i)
			continue
		}

		rows := value.([]interface{})

		if len(rows) != 2 {
			t.Errorf("Test case %d: Got %d rows, expected 2", i, len(rows))
			continue
		}

		// Verify version
		versionRow := rows[0].([]string)

		if len(versionRow) != 2 {
			t.Errorf("Test case %d: Got %d columns, expected 2", i, len(versionRow))
			continue
		}

		if versionRow[0] != "Version number" {
			t.Errorf(`Test case %d: First column was "%s", expected "Version number"`, i, versionRow[0])
		}
		if versionRow[1] != test.snap.Version {
			t.Errorf(`Test case %d: Second column was "%s", expected "%s"`, i, versionRow[1], test.snap.Version)
		}

		// Verify size
		sizeRow := rows[1].([]string)

		if len(sizeRow) != 2 {
			t.Errorf("Test case %d: Got %d columns, expected 2", i, len(sizeRow))
			continue
		}

		if sizeRow[0] != "Size" {
			t.Errorf(`Test case %d: First column was "%s", expected "Size"`, i, sizeRow[0])
		}
		if sizeRow[1] != "124 kB" {
			t.Errorf(`Test case %d: Second column was "%s", expected "124 kB"`, i, sizeRow[1])
		}
	}
}
