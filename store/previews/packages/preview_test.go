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

package packages

import (
	"github.com/snapcore/snapd/client"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/previews/fakes"
	"launchpad.net/unity-scope-snappy/store/previews/packages/templates"
	"reflect"
	"testing"
)

var (
	emptyMetadata     = operation.Metadata{}
	installMetadata   = operation.Metadata{InstallRequested: true, ObjectPath: "/foo/1"}
	uninstallMetadata = operation.Metadata{UninstallConfirmed: true, ObjectPath: "/foo/1"}
)

// Data for both TestNewPreview and TestPreview_generate.
var previewTests = []struct {
	status           string
	metadata         operation.Metadata
	expectedTemplate interface{}
}{
	// No metadata
	{client.StatusInstalled, emptyMetadata, &templates.InstalledTemplate{}},
	{client.StatusAvailable, emptyMetadata, &templates.StoreTemplate{}},
	{client.StatusRemoved, emptyMetadata, &templates.StoreTemplate{}},
	{client.StatusActive, emptyMetadata, &templates.InstalledTemplate{}},

	// Metadata requesting install
	{client.StatusInstalled, installMetadata, &templates.InstalledTemplate{}},
	{client.StatusAvailable, installMetadata, &templates.InstallingTemplate{}},
	{client.StatusRemoved, installMetadata, &templates.InstallingTemplate{}},

	// Metadata requesting uninstall
	{client.StatusInstalled, uninstallMetadata, &templates.UninstallingTemplate{}},
	{client.StatusActive, uninstallMetadata, &templates.UninstallingTemplate{}},
	{client.StatusAvailable, uninstallMetadata, &templates.StoreTemplate{}},
	{client.StatusRemoved, uninstallMetadata, &templates.StoreTemplate{}},
}

// Test typical NewPreview usage.
func TestNewPreview(t *testing.T) {
	for i, test := range previewTests {
		snap := client.Snap{Status: test.status}

		preview, err := NewPreview(snap, nil, test.metadata)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error: %s", i, err)
			continue
		}

		templateType := reflect.TypeOf(preview.template)
		expectedTemplateType := reflect.TypeOf(test.expectedTemplate)
		if templateType != expectedTemplateType {
			t.Errorf(`Test case %d: Template type was "%s", expected "%s"`, i, templateType, expectedTemplateType)
		}
	}
}

// Test typical Generate usage, and verify that it conforms to store design.
func TestPreview_generate(t *testing.T) {
	for i, test := range previewTests {
		preview, err := NewPreview(client.Snap{
			ID:           "package1",
			Name:         "package1",
			Version:      "0.1",
			Developer:    "bar",
			Description:  "baz",
			Icon:         "http://fake",
			Status:       test.status,
			DownloadSize: 123456,
			Type:         "app",
		}, nil, test.metadata)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error while creating package preview: %s", i, err)
			continue
		}

		receiver := new(fakes.FakeWidgetReceiver)

		err = preview.Generate(receiver)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error while generating preview: %s", i, err)
		}

		if len(receiver.Widgets) != 4 {
			// Exit here so we don't index out of bounds later
			t.Fatalf("Test case %d: Got %d widgets, expected 4", i, len(receiver.Widgets))
		}

		widget := receiver.Widgets[0]
		if widget.WidgetType() != "header" {
			t.Errorf("Test case %d: Expected header to be first widget", i)
		}

		widget = receiver.Widgets[1]

		switch test.expectedTemplate.(type) {
		case *templates.InstallingTemplate:
		case *templates.UninstallingTemplate:
			if widget.WidgetType() != "progress" {
				t.Errorf("Test case %d: Expected progress to be second widget", i)
			}
		default:
			if widget.WidgetType() != "actions" {
				t.Errorf("Test case %d: Expected actions to be second widget", i)
			}
		}

		widget = receiver.Widgets[2]
		if widget.WidgetType() != "text" {
			t.Errorf("Test case %d: Expected info to be the third widget", i)
		}

		widget = receiver.Widgets[3]
		if widget.WidgetType() != "table" {
			t.Errorf("Test case %d: Expected updates table to be the fourth widget", i)
		}
	}
}
