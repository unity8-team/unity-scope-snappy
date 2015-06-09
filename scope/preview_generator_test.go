package scope

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"reflect"
	"testing"
)

// Data for TestNewPreview
var newPreviewTests = []struct {
	status       webdm.Status
	scopeData    *ProgressHack
	expectedType string
}{
	{webdm.StatusUndefined, nil, "*scope.StorePreview"},
	{webdm.StatusInstalled, nil, "*scope.InstalledPreview"},
	{webdm.StatusNotInstalled, nil, "*scope.StorePreview"},
	{webdm.StatusInstalling, nil, "*scope.StorePreview"},
	{webdm.StatusUninstalling, nil, "*scope.StorePreview"},

	// Temporary progress-related preview tests
	{webdm.StatusNotInstalled, &ProgressHack{webdm.StatusInstalled}, "*scope.InstallingPreview"},
	{webdm.StatusInstalling, &ProgressHack{webdm.StatusInstalled}, "*scope.InstallingPreview"},
	{webdm.StatusInstalled, &ProgressHack{webdm.StatusInstalled}, "*scope.InstalledPreview"},
	{webdm.StatusInstalled, &ProgressHack{webdm.StatusNotInstalled}, "*scope.UninstallingPreview"},
	{webdm.StatusUninstalling, &ProgressHack{webdm.StatusNotInstalled}, "*scope.UninstallingPreview"},
	{webdm.StatusNotInstalled, &ProgressHack{webdm.StatusNotInstalled}, "*scope.StorePreview"},
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
		if previewType.String() != test.expectedType {
			t.Errorf(`Test case %d: Preview type was "%s", expected "%s"`, i, previewType, test.expectedType)
		}
	}
}

// Test that ProgressHack cannot be retrieved if it's nil
func TestRetrieveProgressHack(t *testing.T) {
	metadata := scopes.NewActionMetadata("us", "phone")
	if retrieveProgressHack(metadata, nil) {
		t.Error("Expected retrieval to fail due to nil ProgressHack")
	}
}
