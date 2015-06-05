package main

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical NewInstalledPreview usage.
func TestNewInstalledPreview(t *testing.T) {
	preview, err := NewInstalledPreview(webdm.Package{
		Id:     "package1",
		Status: webdm.StatusInstalled,
	})
	if err != nil {
		t.Error("Unexpected error creating new installed preview: %s", err)
	}

	if preview.snap.Id != "package1" {
		t.Error(`Preview snap's ID is "%s", expected "package1"`, preview.snap.Id)
	}
}

// Make sure an error occurs if the package is not installed
func TestNewInstalledPreview_notInstalled(t *testing.T) {
	_, err := NewInstalledPreview(webdm.Package{Status: webdm.StatusNotInstalled})
	if err == nil {
		t.Error("Expected an error if the package is not installed")
	}
}

// Test typical Generate usage. This test is perhaps a bit rigid/fragile, but
// it's enforcing the store design.
func TestInstalledPreview_generate(t *testing.T) {
	preview, _ := NewInstalledPreview(
		webdm.Package{
			Id:           "package1",
			Name:         "package1",
			Origin:       "foo",
			Version:      "0.1",
			Vendor:       "bar",
			Description:  "baz",
			IconUrl:      "http://fake",
			Status:       webdm.StatusInstalled,
			DownloadSize: 123456,
			Type:         "oem",
		})

	receiver := new(FakeWidgetReceiver)

	err := preview.Generate(receiver)
	if err != nil {
		t.Errorf("Unexpected error while generating preview: %s", err)
	}

	if len(receiver.widgets) != 4 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Got %d widgets, expected 4", len(receiver.widgets))
	}

	widget := receiver.widgets[0]
	if widget.WidgetType() == "header" {
		verifyInstalledHeaderWidget(t, widget)
	} else {
		t.Error("Expected header to be first widget")
	}

	widget = receiver.widgets[1]
	if widget.WidgetType() == "actions" {
		verifyInstalledActionsWidget(t, widget)
	} else {
		t.Error("Expected actions to be second widget")
	}

	widget = receiver.widgets[2]
	if widget.WidgetType() == "text" {
		verifyInstalledInfoWidget(t, widget, preview.snap.Description)
	} else {
		t.Error("Expected info to be the third widget")
	}

	widget = receiver.widgets[3]
	if widget.WidgetType() == "table" {
		verifyInstalledUpdatesWidget(t, widget, preview.snap.Version)
	} else {
		t.Error("Expected updates table to be the fourth widget")
	}
}

// Test that Generate fails if the package is not installed
func TestInstalledPreview_generate_notInstalled(t *testing.T) {
	preview := InstalledPreview{
		snap: webdm.Package{
			Status: webdm.StatusNotInstalled,
		},
	}

	receiver := new(FakeWidgetReceiver)

	err := preview.Generate(receiver)
	if err == nil {
		t.Error("Expected an error if the package is not installed")
	}
}

// verifyInstalledHeaderWidget is used to verify that a header widget matches
// what we expect.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Header widget to verify.
func verifyInstalledHeaderWidget(t *testing.T, widget scopes.PreviewWidget) {
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

// verifyInstalledActionsWidget is used to verify that an actions widget
// matches what we expect.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Actions widget to verify.
func verifyInstalledActionsWidget(t *testing.T, widget scopes.PreviewWidget) {
	value, ok := widget["actions"]
	if !ok {
		t.Error("Expected actions widget to include actions")
	}

	actionsInterfaces := value.([]interface{})

	if len(actionsInterfaces) != 2 {
		// Exit here so we don't index out of bounds
		t.Fatalf("Actions widget has %d actions, expected 2", len(actionsInterfaces))
	}

	// Verify the open action
	action := actionsInterfaces[0].(map[string]interface{})
	value, ok = action["id"]
	if !ok {
		t.Error("Expected open action to have an id")
	}
	if value != ActionOpen {
		t.Errorf(`Expected open action's ID to be "%d"`, ActionOpen)
	}

	value, ok = action["label"]
	if !ok {
		t.Error("Expected open action to have a label")
	}
	if value != "Open" {
		t.Error(`Expected open action's label to be "Open"`)
	}

	// Verify the uninstall action
	action = actionsInterfaces[1].(map[string]interface{})
	value, ok = action["id"]
	if !ok {
		t.Error("Expected uninstall action to have an id")
	}
	if value != ActionUninstall {
		t.Errorf(`Expected uninstall action's ID to be "%d"`, ActionUninstall)
	}

	value, ok = action["label"]
	if !ok {
		t.Error("Expected uninstall action to have a label")
	}
	if value != "Uninstall" {
		t.Error(`Expected uninstall action's label to be "Uninstall"`)
	}
}

// verifyInstalledInfoWidget is used to verify that a text widget containing
// generic information matches what we expect.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Text widget to verify.
// expectedDescription: Description expected in the text widget.
func verifyInstalledInfoWidget(t *testing.T, widget scopes.PreviewWidget, expectedDescription string) {
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
	if value != expectedDescription {
		t.Errorf(`Got "%s" as the description, expected "%s"`, value, expectedDescription)
	}
}

// verifyInstalledUpdatesWidget is used to verify that a table widget containing
// update information matches what we expect.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Table widget to verify.
// expectedVersion: Version expected in the table widget.
func verifyInstalledUpdatesWidget(t *testing.T, widget scopes.PreviewWidget, expectedVersion string) {
	// Verify title
	value, ok := widget["title"]
	if !ok {
		t.Error("Expected updates table to include a title")
	}
	if value != "Updates" {
		t.Error(`Expected updates table's title to be "Updates"`)
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
		t.Error(`Expected first column to be "Version number"`)
	}
	if versionRow[1] != expectedVersion {
		t.Error(`Expeced second column to be "%s"`, expectedVersion)
	}
}
