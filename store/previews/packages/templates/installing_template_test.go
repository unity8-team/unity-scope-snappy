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
	"github.com/godbus/dbus"
	"github.com/snapcore/snapd/client"
	"testing"
)

// Data for InstallingTemplate tests
var installingTemplateTests = []struct {
	snap        client.Snap
}{
	{client.Snap{ID: "package1", Status: client.StatusAvailable, Version: "0.1", DownloadSize: 123456}},
	{client.Snap{ID: "package1", Status: client.StatusRemoved, Version: "0.1", DownloadSize: 123456}},
}

// Test typical NewInstallingTemplate usage.
func TestNewInstallingTemplate(t *testing.T) {
	for i, test := range installingTemplateTests {
		template, err := NewInstallingTemplate(test.snap, nil, "/foo/1")
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		if template.snap.ID != test.snap.ID {
			t.Errorf(`Test case %d: Template snap's ID is "%s", expected "%s"`, i, template.snap.ID, test.snap.ID)
		}
	}
}

// Test that calling NewInstallingTemplate with an invalid object path results
// in an error.
func TestNewInstallingTemplate_invalidObjectPath(t *testing.T) {
	_, err := NewInstallingTemplate(client.Snap{}, nil, "invalid")
	if err == nil {
		t.Error("Expected an error due to invalid object path")
	}
}

// Test that the actions widget conforms to the store design.
func TestInstallingTemplate_actionsWidget(t *testing.T) {
	for i, test := range installingTemplateTests {
		template, err := NewInstallingTemplate(test.snap, nil, "/foo/1")
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		widget := template.ActionsWidget()

		value, ok := widget["source"]
		if !ok {
			t.Errorf("Test case %d: Expected progress widget to include source", i)
			continue
		}

		// Verify the progress widget
		progressWidget := value.(map[string]interface{})

		value, ok = progressWidget["dbus-name"]
		if !ok {
			t.Errorf("Test case %d: Expected progress widget to have a dbus-name", i)
		}
		if value != "com.canonical.applications.WebdmPackageManager" {
			t.Errorf(`Test case %d: Progress widget's dbus-name was "%s", expected "com.canonical.applications.WebdmPackageManager"`, i, value)
		}

		value, ok = progressWidget["dbus-object"]
		if !ok {
			t.Errorf("Test case %d: Expected progress widget to have a dbus-object", i)
		}
		if value != dbus.ObjectPath("/foo/1") {
			t.Errorf(`Test case %d: Progress widget's dbus-object was "%s", expected "/foo/1"`, i, value)
		}
	}
}
