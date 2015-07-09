package previews

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/previews/packages"
	"launchpad.net/unity-scope-snappy/webdm"
	"reflect"
	"testing"
)

// Data for TestNewPreview.
var newPreviewTests = []struct {
	status    webdm.Status
	scopeData *operation.Metadata
	expected  interface{}
}{
	{webdm.StatusUndefined, nil, &packages.Preview{}},
	{webdm.StatusInstalled, nil, &packages.Preview{}},
	{webdm.StatusNotInstalled, nil, &packages.Preview{}},
	{webdm.StatusInstalling, nil, &packages.Preview{}},
	{webdm.StatusUninstalling, nil, &packages.Preview{}},

	// Uninstallation confirmation test cases
	{webdm.StatusUndefined, &operation.Metadata{UninstallRequested: true}, &ConfirmUninstallPreview{}},
	{webdm.StatusInstalled, &operation.Metadata{UninstallRequested: true}, &ConfirmUninstallPreview{}},
}

// Test typical NewPreview usage.
func TestNewPreview(t *testing.T) {
	for i, test := range newPreviewTests {
		snap := webdm.Package{Status: test.status}
		metadata := scopes.NewActionMetadata("us", "phone")

		metadata.SetScopeData(test.scopeData)

		preview, err := NewPreview(snap, metadata)
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
