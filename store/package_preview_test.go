package store

import (
	"launchpad.net/unity-scope-snappy/webdm"
	"reflect"
	"testing"
)

// Data for both TestNewPackagePreview and TestPackagePreview_generate.
var packagePreviewTests = []struct {
	status           webdm.Status
	expectedTemplate interface{}
}{
	{webdm.StatusUndefined, &StorePackagePreviewTemplate{}},
	{webdm.StatusInstalled, &InstalledPackagePreviewTemplate{}},
	{webdm.StatusNotInstalled, &StorePackagePreviewTemplate{}},
	{webdm.StatusInstalling, &StorePackagePreviewTemplate{}},
	{webdm.StatusUninstalling, &StorePackagePreviewTemplate{}},
}

// Test typical NewPackagePreview usage.
func TestNewPackagePreview(t *testing.T) {
	for i, test := range packagePreviewTests {
		snap := webdm.Package{Status: test.status}

		preview, err := NewPackagePreview(snap)
		if err != nil {
			t.Fatalf("Test case %d: Unexpected error: %s", i, err)
		}

		templateType := reflect.TypeOf(preview.template)
		expectedTemplateType := reflect.TypeOf(test.expectedTemplate)
		if templateType != expectedTemplateType {
			t.Errorf(`Test case %d: Template type was "%s", expected "%s"`, i, templateType, expectedTemplateType)
		}
	}
}

// Test typical Generate usage, and verify that it conforms to store design.
func TestPackagePreview_generate(t *testing.T) {
	for i, test := range packagePreviewTests {
		preview, _ := NewPackagePreview(webdm.Package{
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
		})

		receiver := new(FakeWidgetReceiver)

		err := preview.Generate(receiver)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error while generating preview: %s", i, err)
		}

		if len(receiver.widgets) != 4 {
			// Exit here so we don't index out of bounds later
			t.Fatalf("Test case %d: Got %d widgets, expected 4", i, len(receiver.widgets))
		}

		widget := receiver.widgets[0]
		if widget.WidgetType() != "header" {
			t.Errorf("Test case %d: Expected header to be first widget", i)
		}

		widget = receiver.widgets[1]
		if widget.WidgetType() != "actions" {
			t.Errorf("Test case %d: Expected actions to be second widget", i)
		}

		widget = receiver.widgets[2]
		if widget.WidgetType() != "text" {
			t.Errorf("Test case %d: Expected info to be the third widget", i)
		}

		widget = receiver.widgets[3]
		if widget.WidgetType() != "table" {
			t.Errorf("Test case %d: Expected updates table to be the fourth widget", i)
		}
	}
}
