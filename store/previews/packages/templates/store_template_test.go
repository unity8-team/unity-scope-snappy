package templates

import (
	"launchpad.net/unity-scope-snappy/store/actions"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

// Make sure an error occurs if the package is installed
func TestNewStoreTemplate_installed(t *testing.T) {
	_, err := NewStoreTemplate(webdm.Package{
		Status: webdm.StatusInstalled,
	})

	if err == nil {
		t.Error("Expected an error if the package is installed")
	}
}

// Data for StoreTemplate tests
var storeTemplateTests = []struct {
	snap webdm.Package
}{
	{webdm.Package{Id: "package1", Status: webdm.StatusUndefined, Version: "0.1", DownloadSize: 123456}},
	{webdm.Package{Id: "package1", Status: webdm.StatusNotInstalled, Version: "0.1", DownloadSize: 123456}},
	{webdm.Package{Id: "package1", Status: webdm.StatusInstalling, Version: "0.1", DownloadSize: 123456}},
	{webdm.Package{Id: "package1", Status: webdm.StatusUninstalling, Version: "0.1", DownloadSize: 123456}},
}

// Test typical NewStoreTemplate usage.
func TestNewStoreTemplate(t *testing.T) {
	for i, test := range storeTemplateTests {
		template, err := NewStoreTemplate(test.snap)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		if template.snap.Id != test.snap.Id {
			t.Errorf(`Template snap's ID is "%s", expected "%s"`, template.snap.Id, test.snap.Id)
		}
	}
}

// Test that the header widget conforms to the store design.
func TestStoreTemplate_headerWidget(t *testing.T) {
	for i, test := range storeTemplateTests {
		template, err := NewStoreTemplate(test.snap)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		widget := template.HeaderWidget()

		// Check generic attributes
		value, ok := widget["attributes"]
		if !ok {
			t.Errorf("Test case %d: Expected header attributes to include generic attributes", i)
			continue
		}

		attributes := value.([]interface{})
		if len(attributes) != 1 {
			t.Errorf("Test case %d: Got %d generic attributes for header, expected 1", i, len(attributes))
			continue
		}

		attribute := attributes[0].(map[string]interface{})
		value, ok = attribute["value"]
		if !ok {
			t.Errorf(`Test case %d: Expected generic header attribute to have "value" key`, i)
			continue
		}
		if value != "FREE" {
			t.Errorf(`Test case %d: Generic header attribute was "%s", expected "FREE"`, i, value)
		}
	}
}

// Test that the actions widget conforms to the store design.
func TestStoreTemplate_actionsWidget(t *testing.T) {
	for i, test := range storeTemplateTests {
		template, err := NewStoreTemplate(test.snap)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		widget := template.ActionsWidget()

		value, ok := widget["actions"]
		if !ok {
			t.Errorf("Test case %d: Expected actions widget to include actions", i)
			continue
		}

		actionsInterfaces := value.([]interface{})

		if len(actionsInterfaces) != 1 {
			t.Errorf("Test case %d: Actions widget has %d actions, expected 1", i, len(actionsInterfaces))
			continue
		}

		// Verify the install action
		action := actionsInterfaces[0].(map[string]interface{})
		value, ok = action["id"]
		if !ok {
			t.Errorf("Test case %d: Expected install action to have an id", i)
		}
		if value != actions.ActionInstall {
			t.Errorf(`Test case %d: Install action's ID was "%s", expected "%s"`, i, value, actions.ActionInstall)
		}

		value, ok = action["label"]
		if !ok {
			t.Errorf("Test case %d: Expected install action to have a label", i)
		}
		if value != "Install" {
			t.Errorf(`Test case %d: Install action's label was "%s", expected "Install"`, i, value)
		}
	}
}

// Test that the updates widget conforms to the store design.
func TestStoreTemplate_updatesWidget(t *testing.T) {
	for i, test := range storeTemplateTests {
		template, err := NewStoreTemplate(test.snap)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error creating template: %s", i, err)
			continue
		}

		widget := template.UpdatesWidget()

		// Verify title
		value, ok := widget["title"]
		if !ok {
			t.Errorf("Test case %d: Expected updates table to include a title", i)
		}
		if value != "Updates" {
			t.Errorf(`Test case %d: Updates table's title was "%s", expected "Updates"`, i, value)
		}

		// Verify table rows
		value, ok = widget["values"]
		if !ok {
			t.Errorf("Test case %d: Expected updates table to include values", i)
			continue
		}

		rows := value.([]interface{})

		if len(rows) != 2 {
			t.Errorf("Test case %d: Got %d rows, expected 2", i, len(rows))
			continue
		}

		// Verify version
		versionRow := rows[0].([]string)

		if len(versionRow) != 2 {
			t.Errorf("Test case %d: Got %d columns, expected 2", i, len(versionRow))
			continue
		}

		if versionRow[0] != "Version number" {
			t.Errorf(`Test case %d: First column was "%s", expected "Version number"`, i, versionRow[0])
		}
		if versionRow[1] != test.snap.Version {
			t.Errorf(`Test case %d: Second column was "%s", expected "%s"`, i, versionRow[1], test.snap.Version)
		}

		// Verify size
		sizeRow := rows[1].([]string)

		if len(sizeRow) != 2 {
			t.Errorf("Test case %d: Got %d columns, expected 2", i, len(sizeRow))
			continue
		}

		if sizeRow[0] != "Size" {
			t.Errorf(`Test case %d: First column was "%s", expected "Size"`, i, sizeRow[0])
		}
		if sizeRow[1] != "124 kB" {
			t.Errorf(`Test case %d: Second column was "%s", expected "124 kB"`, i, sizeRow[1])
		}
	}
}
