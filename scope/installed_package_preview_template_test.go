package scope

import (
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical NewInstalledPackagePreviewTemplate usage.
func TestNewInstalledPackagePreviewTemplate(t *testing.T) {
	template, err := NewInstalledPackagePreviewTemplate(webdm.Package{
		Id:     "package1",
		Status: webdm.StatusInstalled,
	})
	if err != nil {
		t.Errorf("Unexpected error creating new installed preview: %s", err)
	}

	if template.snap.Id != "package1" {
		t.Errorf(`Template snap's ID is "%s", expected "package1"`, template.snap.Id)
	}
}

// Make sure an error occurs if the package is not installed
func TestNewInstalledPackagePreviewTemplate_notInstalled(t *testing.T) {
	_, err := NewInstalledPackagePreviewTemplate(webdm.Package{
		Status: webdm.StatusNotInstalled,
	})

	if err == nil {
		t.Error("Expected an error if the package is not installed")
	}
}

// Test that the actions widget conforms to the store design.
func TestNewInstalledPackagePreviewTemplate_actionsWidget(t *testing.T) {
	template, _ := NewInstalledPackagePreviewTemplate(webdm.Package{
		Status: webdm.StatusInstalled,
	})

	widget := template.actionsWidget()

	value, ok := widget["actions"]
	if !ok {
		t.Fatal("Expected actions widget to include actions")
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

// Test that the updates widget conforms to the store design.
func TestNewInstalledPackagePreviewTemplate_updatesWidget(t *testing.T) {
	snap := webdm.Package{
		Version:       "0.1",
		InstalledSize: 123456,
		Status:        webdm.StatusInstalled,
	}
	template, _ := NewInstalledPackagePreviewTemplate(snap)

	widget := template.updatesWidget()

	// Verify title
	value, ok := widget["title"]
	if !ok {
		t.Error("Expected updates table to include a title")
	}
	if value != "Updates" {
		t.Error(`Expected updates table's title to be "Updates"`)
	}

	// Verify table rows
	value, ok = widget["values"]
	if !ok {
		t.Fatal("Expected updates table to include values")
	}

	rows := value.([]interface{})

	if len(rows) != 2 {
		// Exit now so we don't index out of bounds
		t.Fatalf("Got %d rows, expected 2", len(rows))
	}

	// Verify version
	versionRow := rows[0].([]string)

	if len(versionRow) != 2 {
		// Exit now so we don't index out of bounds
		t.Fatalf("Got %d columns, expected 2", len(versionRow))
	}

	if versionRow[0] != "Version number" {
		t.Errorf(`First column was "%s", expected "Version number"`, versionRow[0])
	}
	if versionRow[1] != snap.Version {
		t.Errorf(`Second column was "%s", expected "%s"`, versionRow[1], snap.Version)
	}

	// Verify size
	sizeRow := rows[1].([]string)

	if len(sizeRow) != 2 {
		// Exit now do we don't index out of bounds
		t.Fatalf("Got %d columns, expected 2", len(sizeRow))
	}

	if sizeRow[0] != "Size" {
		t.Error(`First column was "%s", expected "Size"`, sizeRow[0])
	}
	if sizeRow[1] != "124 KB" {
		t.Errorf(`Second column was "%s", expected "124 KB"`, sizeRow[1])
	}
}
