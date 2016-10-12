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

// Data for UninstallingTemplate tests
var uninstallingTemplateTests = []struct {
	snap        client.Snap
}{
	{client.Snap{ID: "package1", Status: client.StatusInstalled, Version: "0.1", InstalledSize: 123456}},
	{client.Snap{ID: "package1", Status: client.StatusActive, Version: "0.1", DownloadSize: 123456}},
}

// Test typical NewUninstallingTemplate usage.
func TestNewUninstallingTemplate(t *testing.T) {
	for i, test := range uninstallingTemplateTests {
		template, err := NewUninstallingTemplate(test.snap, "/foo/1")
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		if template.snap.ID != test.snap.ID {
			t.Errorf(`Test case %d: Template snap's ID is "%s", expected "%s"`, i, template.snap.ID, test.snap.ID)
		}
	}
}

// Test that calling NewUninstallingTemplate with an invalid object path results
// in an error.
func TestNewUninstallingTemplate_invalidObjectPath(t *testing.T) {
	_, err := NewUninstallingTemplate(client.Snap{}, "invalid")
	if err == nil {
		t.Error("Expected an error due to invalid object path")
	}
}

// Test that the actions widget conforms to the store design.
func TestUninstallingTemplate_actionsWidget(t *testing.T) {
	for i, test := range uninstallingTemplateTests {
		template, err := NewUninstallingTemplate(test.snap, "/foo/1")
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
