package packages

import (
	"launchpad.net/unity-scope-snappy/store/previews/fakes"
	"launchpad.net/unity-scope-snappy/store/previews/packages/templates"
	"launchpad.net/unity-scope-snappy/webdm"
	"reflect"
	"testing"
)

// Data for both TestNewPreview and TestPreview_generate.
var previewTests = []struct {
	status           webdm.Status
	expectedTemplate interface{}
}{
	{webdm.StatusUndefined, &templates.StoreTemplate{}},
	{webdm.StatusInstalled, &templates.InstalledTemplate{}},
	{webdm.StatusNotInstalled, &templates.StoreTemplate{}},
	{webdm.StatusInstalling, &templates.StoreTemplate{}},
	{webdm.StatusUninstalling, &templates.StoreTemplate{}},
}

// Test typical NewPreview usage.
func TestNewPreview(t *testing.T) {
	for i, test := range previewTests {
		snap := webdm.Package{Status: test.status}

		preview, err := NewPreview(snap)
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
func TestPreview_generate(t *testing.T) {
	for i, test := range previewTests {
		preview, _ := NewPreview(webdm.Package{
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

		receiver := new(fakes.FakeWidgetReceiver)

		err := preview.Generate(receiver)
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
		if widget.WidgetType() != "actions" {
			t.Errorf("Test case %d: Expected actions to be second widget", i)
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
