package packages

import (
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/previews/interfaces"
	"launchpad.net/unity-scope-snappy/store/previews/packages/templates"
	"launchpad.net/unity-scope-snappy/webdm"
	"reflect"
	"testing"
)

var (
	emptyMetadata     = operation.Metadata{}
	installMetadata   = operation.Metadata{InstallRequested: true, ObjectPath: "/foo/1"}
	uninstallMetadata = operation.Metadata{UninstallConfirmed: true, ObjectPath: "/foo/1"}
)

// Data for both TestNewPreview_invalidMetadata
var invalidMetadataTests = []struct {
	status   webdm.Status
	metadata operation.Metadata
}{
	{webdm.StatusUninstalling, installMetadata},
	{webdm.StatusInstalling, uninstallMetadata},
	{webdm.StatusUndefined, uninstallMetadata},
}

// Test that calling NewPreview with invalid metadata results in an error.
func TestNewPreview_invalidMetadata(t *testing.T) {
	for i, test := range invalidMetadataTests {
		snap := webdm.Package{Status: test.status}

		_, err := NewPreview(snap, test.metadata)
		if err == nil {
			t.Errorf("Test case %d: Expected an error due to invalid metadata", i)
		}
	}
}

// Data for both TestNewPreview and TestPreview_generate.
var previewTests = []struct {
	status           webdm.Status
	metadata         operation.Metadata
	expectedTemplate interface{}
}{
	// No metadata
	{webdm.StatusUndefined, emptyMetadata, &templates.StoreTemplate{}},
	{webdm.StatusInstalled, emptyMetadata, &templates.InstalledTemplate{}},
	{webdm.StatusNotInstalled, emptyMetadata, &templates.StoreTemplate{}},
	{webdm.StatusInstalling, emptyMetadata, &templates.StoreTemplate{}},
	{webdm.StatusUninstalling, emptyMetadata, &templates.StoreTemplate{}},

	// Metadata requesting install
	{webdm.StatusUndefined, installMetadata, &templates.InstallingTemplate{}},
	{webdm.StatusInstalled, installMetadata, &templates.InstalledTemplate{}},
	{webdm.StatusNotInstalled, installMetadata, &templates.InstallingTemplate{}},
	{webdm.StatusInstalling, installMetadata, &templates.InstallingTemplate{}},

	// Metadata requesting uninstall
	{webdm.StatusInstalled, uninstallMetadata, &templates.UninstallingTemplate{}},
	{webdm.StatusNotInstalled, uninstallMetadata, &templates.StoreTemplate{}},
	{webdm.StatusUninstalling, uninstallMetadata, &templates.UninstallingTemplate{}},
}

// Test typical NewPreview usage.
func TestNewPreview(t *testing.T) {
	for i, test := range previewTests {
		snap := webdm.Package{Status: test.status}

		preview, err := NewPreview(snap, test.metadata)
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
		preview, err := NewPreview(webdm.Package{
			Id:           "package1",
			Name:         "package1",
			Origin:       "foo",
			Version:      "0.1",
			Vendor:       "bar",
			Description:  "baz",
			IconUrl:      "http://fake",
			Status:       test.status,
			DownloadSize: 123456,
			Type:         "oem",
		}, test.metadata)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error while creating package preview: %s", i, err)
			continue
		}

		receiver := new(interfaces.FakeWidgetReceiver)

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
