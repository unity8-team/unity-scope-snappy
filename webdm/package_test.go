package webdm

import (
	"testing"
)

// Data for TestMarshalJSON
var marshalJSONTests = []struct {
	status      Status
	expected    string
	shouldError bool
}{
	{StatusUndefined, "", true},
	{StatusNotInstalled, `"uninstalled"`, false},
	{StatusInstalled, `"installed"`, false},
	{StatusInstalling, `"installing"`, false},
	{StatusUninstalling, `"uninstalling"`, false},
}

// Test typical MarshalJSON usage
func TestMarshalJSON(t *testing.T) {
	for i, test := range marshalJSONTests {
		json, err := test.status.MarshalJSON()

		if test.shouldError && err == nil {
			t.Errorf("Test case %d: Expected an error", i)
		}

		if !test.shouldError && err != nil {
			t.Errorf("Test case %d: Unexpected error: %s", i, err)
		}

		if !test.shouldError && string(json) != test.expected {
			t.Errorf("Test case %d: Got %s, expected %s", i, string(json), test.expected)
		}
	}
}

// Data for TestUnmarshalJSON
var unmarshalJSONTests = []struct {
	json        string
	expected    Status
	shouldError bool
}{
	{`"uninstalled"`, StatusNotInstalled, false},
	{`"installed"`, StatusInstalled, false},
	{`"installing"`, StatusInstalling, false},
	{`"uninstalling"`, StatusUninstalling, false},
	{`"foo"`, StatusUndefined, true},
}

// Test typical UnmarshalJSON usage
func TestUnmarshalJSON(t *testing.T) {
	for i, test := range unmarshalJSONTests {
		var status Status
		err := status.UnmarshalJSON([]byte(test.json))
		if status != test.expected {
			t.Errorf("Test case %d: Status was %d, expected %d", i, status, test.expected)
		}

		if test.shouldError && err == nil {
			t.Errorf("Test case %d: Expected an error", i)
		}

		if !test.shouldError && err != nil {
			t.Errorf("Test case %d: Unexpected error: %s", i, err)
		}
	}
}

// Test that UnmarshalJSON dies when called on nil Status
func TestUnmarshalJSON_nilStatus(t *testing.T) {
	var nilStatus *Status = nil
	err := nilStatus.UnmarshalJSON([]byte("this is fake json"))
	if err == nil {
		t.Error("Expected error when Status is nil")
	}
}

// Test typical Installed usage.
func TestPackage_installed(t *testing.T) {
	snap := Package{Status: StatusInstalled}
	if !snap.Installed() {
		t.Error("Expected package to be installed")
	}
}

// Test typical Installing usage.
func TestPackage_installing(t *testing.T) {
	snap := Package{Status: StatusInstalling}
	if !snap.Installing() {
		t.Error("Expected package to be installing")
	}
}

// Test typical NotInstalled usage.
func TestPackage_notInstalled(t *testing.T) {
	snap := Package{Status: StatusNotInstalled}
	if !snap.NotInstalled() {
		t.Error("Expected package to not be installed")
	}
}

// Test typical Uninstalling usage.
func TestPackage_uninstalling(t *testing.T) {
	snap := Package{Status: StatusUninstalling}
	if !snap.Uninstalling() {
		t.Error("Expected package to be uninstalling")
	}
}
