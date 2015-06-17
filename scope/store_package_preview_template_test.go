package scope

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
