package scope

import (
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

var (
	webdmPackage *webdm.Package
	template     *GenericPackagePreviewTemplate
)

func setup() {
	webdmPackage = &webdm.Package{
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
	}

	template = NewGenericPackagePreviewTemplate(*webdmPackage)
}

// Test typical NewGenericPackagePreviewTemplate usage.
func TestNewGenericPackagePreviewTemplate(t *testing.T) {
	setup()

	if template.snap.Id != "package1" {
		t.Errorf(`Template snap's ID was "%s", expected "package1"`, template.snap.Id)
	}
}

// Test that the header widget conforms to the store design.
func TestGenericPackagePreviewTemplate_headerWidget(t *testing.T) {
	setup()

	widget := template.headerWidget()

	if widget.WidgetType() != "header" {
		t.Fatal(`Expected widget type to be "header"`)
	}

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

// Test that the actions widget doesn't actually contain anything, as a generic
// package doesn't have enough information to fill it out.
func TestGenericPackagePreviewTemplate_actionsWidget(t *testing.T) {
	setup()

	widget := template.actionsWidget()

	if widget.WidgetType() != "actions" {
		t.Fatal(`Expected widget type to be "actions"`)
	}
}

// Test that the header widget conforms to the store design.
func TestGenericPackagePreviewTemplate_infoWidget(t *testing.T) {
	setup()

	widget := template.infoWidget()

	if widget.WidgetType() != "text" {
		t.Fatal(`Expected widget type to be "text"`)
	}

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
	if value != webdmPackage.Description {
		t.Errorf(`Got "%s" as the description, expected "%s"`, value, webdmPackage.Description)
	}
}

// Test that the updates widget conforms to the store design.
func TestGenericPackagePreviewTemplate_updatesWidget(t *testing.T) {
	setup()

	widget := template.updatesWidget()

	if widget.WidgetType() != "table" {
		t.Fatal(`Expected widget type to be "table"`)
	}

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
	if versionRow[1] != webdmPackage.Version {
		t.Error(`Expeced second column to be "%s"`, webdmPackage.Version)
	}
}