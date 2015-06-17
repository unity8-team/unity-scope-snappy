package scope

import (
	"launchpad.net/unity-scope-snappy/webdm"
	"reflect"
	"testing"
)

// Data for TestNewPreview
var newPreviewTests = []struct {
	status   webdm.Status
	expected interface{}
}{
	{webdm.StatusUndefined, &PackagePreview{}},
	{webdm.StatusInstalled, &PackagePreview{}},
	{webdm.StatusNotInstalled, &PackagePreview{}},
	{webdm.StatusInstalling, &PackagePreview{}},
	{webdm.StatusUninstalling, &PackagePreview{}},
}

// Test typical NewPreview usage.
func TestNewPreview(t *testing.T) {
	for i, test := range newPreviewTests {
		snap := webdm.Package{Status: test.status}

		preview, err := NewPreview(snap)
		if err != nil {
			t.Errorf("Test case %d: Unexpected error: %s", i, err)
		}

		previewType := reflect.TypeOf(preview)
		expectedType := reflect.TypeOf(test.expected)
		if previewType != expectedType {
			t.Errorf(`Test case %d: Preview type was "%s", expected "%s"`, i, previewType, expectedType)
		}
	}
}
