package templates

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Data for UninstallingTemplate tests
var uninstallingTemplateTests = []struct {
	snap        webdm.Package
	expectError bool
}{
	{webdm.Package{Id: "package1", Status: webdm.StatusUndefined, Version: "0.1", DownloadSize: 123456}, true},
	{webdm.Package{Id: "package1", Status: webdm.StatusInstalled, Version: "0.1", InstalledSize: 123456}, false},
	{webdm.Package{Id: "package1", Status: webdm.StatusNotInstalled, Version: "0.1", DownloadSize: 123456}, true},
	{webdm.Package{Id: "package1", Status: webdm.StatusInstalling, Version: "0.1", DownloadSize: 123456}, true},
	{webdm.Package{Id: "package1", Status: webdm.StatusUninstalling, Version: "0.1", DownloadSize: 123456}, false},
}

// Test typical NewUninstallingTemplate usage.
func TestNewUninstallingTemplate(t *testing.T) {
	for i, test := range uninstallingTemplateTests {
		template, err := NewUninstallingTemplate(test.snap, "/foo/1")
		if err == nil && test.expectError {
			t.Errorf("Test case %d: Expected error due to incorrect status", i)
		} else if err != nil {
			if !test.expectError {
				t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			}
			continue
		}

		if template.snap.Id != test.snap.Id {
			t.Errorf(`Test case %d: Template snap's ID is "%s", expected "%s"`, i, template.snap.Id, test.snap.Id)
		}
	}
}

// Test that calling NewUninstallingTemplate with an invalid object path results
// in an error.
func TestNewUninstallingTemplate_invalidObjectPath(t *testing.T) {
	_, err := NewUninstallingTemplate(webdm.Package{}, "invalid")
	if err == nil {
		t.Error("Expected an error due to invalid object path")
	}
}

// Test that the actions widget conforms to the store design.
func TestUninstallingTemplate_actionsWidget(t *testing.T) {
	for i, test := range uninstallingTemplateTests {
		template, err := NewUninstallingTemplate(test.snap, "/foo/1")
		if err == nil && test.expectError {
			t.Errorf("Test case %d: Expected error due to incorrect status", i)
		} else if err != nil {
			if !test.expectError {
				t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			}
			continue
		}

		widget := template.ActionsWidget()

		value, ok := widget["source"]
		if !ok {
			t.Errorf("Test case %d: Expected progress widget to include source", i)
			continue
		}

		// Verify the progress widget
		progressWidget := value.(map[string]interface{})

		value, ok = progressWidget["dbus-name"]
		if !ok {
			t.Errorf("Test case %d: Expected progress widget to have a dbus-name", i)
		}
		if value != "com.canonical.applications.WebdmPackageManager" {
			t.Errorf(`Test case %d: Progress widget's dbus-name was "%s", expected "com.canonical.applications.WebdmPackageManager"`, i, value)
		}

		value, ok = progressWidget["dbus-object"]
		if !ok {
			t.Errorf("Test case %d: Expected progress widget to have a dbus-object", i)
		}
		if value != dbus.ObjectPath("/foo/1") {
			t.Errorf(`Test case %d: Progress widget's dbus-object was "%s", expected "/foo/1"`, i, value)
		}
	}
}
