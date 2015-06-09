package scope

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical NewUninstallingPreview usage.
func TestNewUninstallingPreview(t *testing.T) {
	preview, err := NewUninstallingPreview(webdm.Package{
		Id:     "package1",
		Status: webdm.StatusInstalled,
	})
	if err != nil {
		t.Error("Unexpected error creating new uninstalling preview: %s", err)
	}

	if preview.snap.Id != "package1" {
		t.Error(`Preview snap's ID is "%s", expected "package1"`, preview.snap.Id)
	}
}

// Make sure an error occurs if the package isn't installed
func TestNewUninstallingPreview_notInstalled(t *testing.T) {
	_, err := NewUninstallingPreview(webdm.Package{Status: webdm.StatusNotInstalled})
	if err == nil {
		t.Error("Expected an error if the package isn't installed")
	}
}

// Test typical Generate usage.
func TestUninstallingPreview_generate(t *testing.T) {
	preview, _ := NewUninstallingPreview(
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

	if len(receiver.widgets) != 2 {
		// Exit here so we don't index out of bounds later
		t.Fatalf("Got %d widgets, expected 2", len(receiver.widgets))
	}

	widget := receiver.widgets[0]
	if widget.WidgetType() == "text" {
		verifyUninstallingTextWidget(t, widget)
	} else {
		t.Error("Expected text to be first widget")
	}

	widget = receiver.widgets[1]
	if widget.WidgetType() == "actions" {
		verifyUninstallingActionsWidget(t, widget)
	} else {
		t.Error("Expected actions to be second widget")
	}
}

// Test that errors are displayed if the snap doesn't successfully uninstall
func TestUninstallingPreview_generate_uninstallFailed(t *testing.T) {
	preview, _ := NewUninstallingPreview(
		webdm.Package{
			Id:      "package1",
			Name:    "package1",
			Status:  webdm.StatusInstalled,
			Message: "Unable to uninstall", // This indicates failure
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
		verifyUninstallingFailedTextWidget(t, widget)
	} else {
		t.Error("Expected text to be first widget")
	}

	widget = receiver.widgets[1]
	if widget.WidgetType() == "actions" {
		verifyUninstallingFailedActionsWidget(t, widget)
	} else {
		t.Error("Expected actions to be second widget")
	}
}

// Test that Generate fails if the package isn't installed
func TestUninstallingPreview_generate_notInstalled(t *testing.T) {
	preview := UninstallingPreview{
		snap: webdm.Package{
			Status: webdm.StatusNotInstalled,
		},
	}

	receiver := new(FakeWidgetReceiver)

	err := preview.Generate(receiver)
	if err == nil {
		t.Error("Expected an error if the package isn't installed")
	}
}

// verifyUninstallingTextWidget is used to verify that a text widget matches
// what we expect.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Text widget to verify.
func verifyUninstallingTextWidget(t *testing.T, widget scopes.PreviewWidget) {
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

// verifyUninstallingFailedTextWidget is used to verify that a text widget
// matches what we expect when the uninstall fails.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Text widget to verify.
func verifyUninstallingFailedTextWidget(t *testing.T, widget scopes.PreviewWidget) {
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

// verifyUninstallingActionsWidget is used to verify that an actions widget
// matches what we expect.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Actions widget to verify.
func verifyUninstallingActionsWidget(t *testing.T, widget scopes.PreviewWidget) {
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
	if value != ActionRefreshUninstalling {
		t.Errorf(`Expected refresh action's ID to be "%d"`, ActionRefreshUninstalling)
	}

	value, ok = action["label"]
	if !ok {
		t.Error("Expected refresh action to have a label")
	}
	if value != "Refresh" {
		t.Error(`Expected refresh action's label to be "Refresh"`)
	}
}

// verifyUninstallingFailedActionsWidget is used to verify that an actions
// widget matches what we expect when an uninstall failed.
//
// Parameters:
// t: Testing handle for when errors occur.
// widget: Actions widget to verify.
func verifyUninstallingFailedActionsWidget(t *testing.T, widget scopes.PreviewWidget) {
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
