package scope

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical NewInstallingPreview usage.
func TestNewInstallingPreview(t *testing.T) {
	preview, err := NewInstallingPreview(webdm.Package{
		Id:     "package1",
		Status: webdm.StatusNotInstalled,
	})
	if err != nil {
		t.Error("Unexpected error creating new installing preview: %s", err)
	}

	if preview.snap.Id != "package1" {
		t.Error(`Preview snap's ID is "%s", expected "package1"`, preview.snap.Id)
	}
}

// Make sure an error occurs if the package is already installed
func TestNewInstallingPreview_installed(t *testing.T) {
	_, err := NewInstallingPreview(webdm.Package{Status: webdm.StatusInstalled})
	if err == nil {
		t.Error("Expected an error if the package is already installed")
	}
}

// Test typical Generate usage.
func TestInstallingPreview_generate(t *testing.T) {
	preview, _ := NewInstallingPreview(
		webdm.Package{
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
		})

	receiver := new(FakeWidgetReceiver)

	err := preview.Generate(receiver)
	if err != nil {
		t.Errorf("Unexpected error while generating preview: %s", err)
	}

	if len(receiver.widgets) != 2 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Got %d widgets, expected 2", len(receiver.widgets))
	}

	widget := receiver.widgets[0]
	if widget.WidgetType() == "text" {
		verifyInstallingTextWidget(t, widget)
	} else {
		t.Error("Expected text to be first widget")
	}

	widget = receiver.widgets[1]
	if widget.WidgetType() == "actions" {
		verifyInstallingActionsWidget(t, widget)
	} else {
		t.Error("Expected actions to be second widget")
	}
}

// Test that errors are displayed if the snap doesn't successfully install
func TestInstallingPreview_generate_installFailed(t *testing.T) {
	preview, _ := NewInstallingPreview(
		webdm.Package{
			Id:      "package1",
			Name:    "package1",
			Status:  webdm.StatusNotInstalled,
			Message: "Unable to install", // This indicates failure
		})

	receiver := new(FakeWidgetReceiver)

	err := preview.Generate(receiver)
	if err != nil {
		t.Errorf("Unexpected error while generating preview: %s", err)
	}

	if len(receiver.widgets) != 2 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Got %d widgets, expected 2", len(receiver.widgets))
	}

	widget := receiver.widgets[0]
	if widget.WidgetType() == "text" {
		verifyInstallingFailedTextWidget(t, widget)
	} else {
		t.Error("Expected text to be first widget")
	}

	widget = receiver.widgets[1]
	if widget.WidgetType() == "actions" {
		verifyInstallingFailedActionsWidget(t, widget)
	} else {
		t.Error("Expected actions to be second widget")
	}
}

// Test that Generate fails if the package is already installed
func TestInstallingPreview_generate_installed(t *testing.T) {
	preview := InstallingPreview{
		snap: webdm.Package{
			Status: webdm.StatusInstalled,
		},
	}

	receiver := new(FakeWidgetReceiver)

	err := preview.Generate(receiver)
	if err == nil {
		t.Error("Expected an error if the package is already installed")
	}
}

// verifyInstallingTextWidget is used to verify that a text widget matches what
// we expect while a snap is installing.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Text widget to verify.
func verifyInstallingTextWidget(t *testing.T, widget scopes.PreviewWidget) {
	// Verify title presence
	_, ok := widget["title"]
	if !ok {
		// Exit here so we don't dereference nil
		t.Fatal("Expected text to include a title")
	}

	// Verify text presence
	_, ok = widget["text"]
	if !ok {
		// Exit here so we don't dereference nil
		t.Fatal("Expected text to include actual text")
	}
}

// verifyInstallingFailedTextWidget is used to verify that a text widget
// matches what we expect when the install fails.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Text widget to verify.
func verifyInstallingFailedTextWidget(t *testing.T, widget scopes.PreviewWidget) {
	// Verify title presence
	_, ok := widget["title"]
	if !ok {
		// Exit here so we don't dereference nil
		t.Fatal("Expected text to include a title")
	}

	// Verify text presence
	_, ok = widget["text"]
	if !ok {
		// Exit here so we don't dereference nil
		t.Fatal("Expected text to include actual text")
	}
}

// verifyInstallingActionsWidget is used to verify that an actions widget
// matches what we expect while a snap is installing.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Actions widget to verify.
func verifyInstallingActionsWidget(t *testing.T, widget scopes.PreviewWidget) {
	value, ok := widget["actions"]
	if !ok {
		t.Error("Expected actions widget to include actions")
	}

	actionsInterfaces := value.([]interface{})

	if len(actionsInterfaces) != 1 {
		// Exit here so we don't index out of bounds
		t.Fatalf("Actions widget has %d actions, expected 1", len(actionsInterfaces))
	}

	// Verify the refresh action
	action := actionsInterfaces[0].(map[string]interface{})
	value, ok = action["id"]
	if !ok {
		t.Error("Expected refresh action to have an id")
	}
	if value != ActionRefreshInstalling {
		t.Errorf(`Expected refresh action's ID to be "%d"`, ActionRefreshInstalling)
	}

	value, ok = action["label"]
	if !ok {
		t.Error("Expected refresh action to have a label")
	}
	if value != "Refresh" {
		t.Error(`Expected refresh action's label to be "Refresh"`)
	}
}

// verifyInstallingFailedActionsWidget is used to verify that an actions
// widget matches what we expect when an install failed.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Actions widget to verify.
func verifyInstallingFailedActionsWidget(t *testing.T, widget scopes.PreviewWidget) {
	value, ok := widget["actions"]
	if !ok {
		t.Error("Expected actions widget to include actions")
	}

	actionsInterfaces := value.([]interface{})

	if len(actionsInterfaces) != 1 {
		// Exit here so we don't index out of bounds
		t.Fatalf("Actions widget has %d actions, expected 1", len(actionsInterfaces))
	}

	// Verify the okay action
	action := actionsInterfaces[0].(map[string]interface{})
	value, ok = action["id"]
	if !ok {
		t.Error("Expected okay action to have an id")
	}
	if value != ActionOk {
		t.Errorf(`Expected okay action's ID to be "%d"`, ActionOk)
	}

	value, ok = action["label"]
	if !ok {
		t.Error("Expected okay action to have a label")
	}
	if value != "Okay" {
		t.Error(`Expected okay action's label to be "Okay"`)
	}
}
