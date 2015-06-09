package scope

import (
	"launchpad.net/unity-scope-snappy/webdm"
	"reflect"
	"testing"
)

// Data for TestNewPreview
var newPreviewTests = []struct {
	status       webdm.Status
	expectedType string
}{
	{webdm.StatusUndefined, "*scope.StorePreview"},
	{webdm.StatusInstalled, "*scope.InstalledPreview"},
	{webdm.StatusNotInstalled, "*scope.StorePreview"},
	{webdm.StatusInstalling, "*scope.StorePreview"},
	{webdm.StatusUninstalling, "*scope.StorePreview"},
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
		if previewType.String() != test.expectedType {
			t.Errorf(`Test case %d: Preview type was "%s", expected "%s"`, i, previewType, test.expectedType)
		}
	}
}
