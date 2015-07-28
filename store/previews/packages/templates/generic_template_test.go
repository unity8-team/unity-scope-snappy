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
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

var (
	webdmPackage *webdm.Package
	template     *GenericTemplate
)

func setup() {
	webdmPackage = &webdm.Package{
		Id:           "package1",
		Name:         "package1",
		Origin:       "foo",
		Version:      "0.1",
		Vendor:       "bar",
		Description:  "baz",
		IconUrl:      "http://fake",
		Status:       webdm.StatusNotInstalled,
		DownloadSize: 123456,
		Type:         "oem",
	}

	template = NewGenericTemplate(*webdmPackage)
}

// Test typical NewGenericTemplate usage.
func TestNewGenericTemplate(t *testing.T) {
	setup()

	if template.snap.Id != "package1" {
		t.Errorf(`Template snap's ID was "%s", expected "package1"`, template.snap.Id)
	}
}

// Test that the header widget conforms to the store design.
func TestNewGenericTemplate_headerWidget(t *testing.T) {
	setup()

	widget := template.HeaderWidget()

	if widget.WidgetType() != "header" {
		t.Fatal(`Expected widget type to be "header"`)
	}

	// Grab attribute mappings
	value, ok := widget["components"]
	if !ok {
		// Exit here so we don't index into a nil `components`
		t.Fatal("Expected header to include attribute mappings")
	}

	components := value.(map[string]interface{})

	// Check "title" component
	value, ok = components["title"]
	if !ok {
		t.Error("Expected header attributes to include a title")
	}
	if value != "title" {
		t.Error(`Expected header title attribute to be mapped to "title"`)
	}

	// Check "subtitle" component
	value, ok = components["subtitle"]
	if !ok {
		t.Error("Expected header attributes to include a subtitle")
	}
	if value != "subtitle" {
		t.Error(`Expected header subtitle attribute to be mapped to "subtitle"`)
	}

	// Check mascot attribute
	value, ok = components["mascot"]
	if !ok {
		t.Error("Expected header attributes to include a mascot")
	}
	if value != "art" {
		t.Error(`Expected header mascot attribute to be mapped to "art"`)
	}
}

// Test that the actions widget doesn't actually contain anything, as a generic
// package doesn't have enough information to fill it out.
func TestNewGenericTemplate_actionsWidget(t *testing.T) {
	setup()

	widget := template.ActionsWidget()

	if widget.WidgetType() != "actions" {
		t.Fatal(`Expected widget type to be "actions"`)
	}
}

// Test that the header widget conforms to the store design.
func TestNewGenericTemplate_infoWidget(t *testing.T) {
	setup()

	widget := template.InfoWidget()

	if widget.WidgetType() != "text" {
		t.Fatal(`Expected widget type to be "text"`)
	}

	// Verify title
	value, ok := widget["title"]
	if !ok {
		t.Error("Expected info widget to include a title")
	}
	if value != "Info" {
		t.Error(`Expected info widget's title to be "Info"`)
	}

	// Verify description
	value, ok = widget["text"]
	if !ok {
		t.Error("Expected info widget to include a description")
	}
	if value != webdmPackage.Description {
		t.Errorf(`Got "%s" as the description, expected "%s"`, value, webdmPackage.Description)
	}
}

// Test that the updates widget conforms to the store design.
func TestNewGenericTemplate_updatesWidget(t *testing.T) {
	setup()

	widget := template.UpdatesWidget()

	if widget.WidgetType() != "table" {
		t.Fatal(`Widget type was "%s", expected "table"`, widget.WidgetType())
	}

	// Verify title
	value, ok := widget["title"]
	if !ok {
		t.Error("Expected updates table to include a title")
	}
	if value != "Updates" {
		t.Error(`Updates table's title was "%s", expected "Updates"`, value)
	}

	// Verify version
	value, ok = widget["values"]
	if !ok {
		t.Error("Expected updates table to include values")
	}

	rows := value.([]interface{})

	if len(rows) != 1 {
		// Exit now so we don't index out of bounds
		t.Fatalf("Got %d rows, expected 1", len(rows))
	}

	versionRow := rows[0].([]string)

	if len(versionRow) != 2 {
		// Exit now so we don't index out of bounds
		t.Fatalf("Got %d columns, expected 2", len(versionRow))
	}

	if versionRow[0] != "Version number" {
		t.Errorf(`First column was "%s", expected "Version number"`, versionRow[0])
	}
	if versionRow[1] != webdmPackage.Version {
		t.Errorf(`Second column was "%s", expected "%s"`, versionRow[1], webdmPackage.Version)
	}
}
