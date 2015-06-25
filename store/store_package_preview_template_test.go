package store

import (
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Test typical NewStorePackagePreviewTemplate usage.
func TestNewStorePackagePreviewTemplate(t *testing.T) {
	template, err := NewStorePackagePreviewTemplate(webdm.Package{
		Id:     "package1",
		Status: webdm.StatusNotInstalled,
	})
	if err != nil {
		t.Errorf("Unexpected error creating new store preview: %s", err)
	}

	if template.snap.Id != "package1" {
		t.Errorf(`Template snap's ID is "%s", expected "package1"`, template.snap.Id)
	}
}

// Make sure an error occurs if the package is installed
func TestNewStorePackagePreviewTemplate_installed(t *testing.T) {
	_, err := NewStorePackagePreviewTemplate(webdm.Package{
		Status: webdm.StatusInstalled,
	})

	if err == nil {
		t.Error("Expected an error if the package is installed")
	}
}

// Test that the header widget conforms to the store design.
func TestNewStorePackagePreviewTemplate_headerWidget(t *testing.T) {
	template, _ := NewStorePackagePreviewTemplate(webdm.Package{
		Status: webdm.StatusNotInstalled,
	})

	widget := template.headerWidget()

	// Check generic attributes
	value, ok := widget["attributes"]
	if !ok {
		t.Fatal("Expected header attributes to include generic attributes")
	}

	attributes := value.([]interface{})
	if len(attributes) != 1 {
		// Exit here so we don't index out of bounds
		t.Fatalf("Got %d generic attributes for header, expected 1", len(attributes))
	}

	attribute := attributes[0].(map[string]interface{})
	value, ok = attribute["value"]
	if !ok {
		t.Error(`Expected generic header attribute to have "value" key`)
	}
	if value != "FREE" {
		t.Error(`Expected generic header attribute "value" to be "FREE"`)
	}
}

// Test that the actions widget conforms to the store design.
func TestNewStorePackagePreviewTemplate_actionsWidget(t *testing.T) {
	template, _ := NewStorePackagePreviewTemplate(webdm.Package{
		Status: webdm.StatusNotInstalled,
	})

	widget := template.actionsWidget()

	value, ok := widget["actions"]
	if !ok {
		t.Fatal("Expected actions widget to include actions")
	}

	actionsInterfaces := value.([]interface{})

	if len(actionsInterfaces) != 1 {
		// Exit here so we don't index out of bounds
		t.Fatalf("Actions widget has %d actions, expected 1", len(actionsInterfaces))
	}

	// Verify the install action
	action := actionsInterfaces[0].(map[string]interface{})
	value, ok = action["id"]
	if !ok {
		t.Error("Expected install action to have an id")
	}
	if value != ActionInstall {
		t.Errorf(`Expected install action's ID to be "%d"`, ActionInstall)
	}

	value, ok = action["label"]
	if !ok {
		t.Error("Expected install action to have a label")
	}
	if value != "Install" {
		t.Error(`Expected install action's label to be "Install"`)
	}
}

// Test that the updates widget conforms to the store design.
func TestNewStorePackagePreviewTemplate_updatesWidget(t *testing.T) {
	snap := webdm.Package{
		Version:      "0.1",
		DownloadSize: 123456,
		Status:       webdm.StatusNotInstalled,
	}
	template, _ := NewStorePackagePreviewTemplate(snap)

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
		t.Error(`Expected first column to be "Version number"`)
	}
	if versionRow[1] != snap.Version {
		t.Error(`Expeced second column to be "%s"`, snap.Version)
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
	if sizeRow[1] != "124 kB" {
		t.Errorf(`Second column was "%s", expected "124 kB"`, sizeRow[1])
	}
}
