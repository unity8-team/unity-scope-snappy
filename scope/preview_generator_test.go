package scope

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
	"reflect"
	"testing"
)

// Data for TestNewPreview.
var newPreviewTests = []struct {
	status    webdm.Status
	scopeData *ProgressHack
	expected  interface{}
}{
	{webdm.StatusUndefined, nil, &StorePreview{}},
	{webdm.StatusInstalled, nil, &InstalledPreview{}},
	{webdm.StatusNotInstalled, nil, &StorePreview{}},
	{webdm.StatusInstalling, nil, &StorePreview{}},
	{webdm.StatusUninstalling, nil, &StorePreview{}},
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

// Data for TestNewPreview_progressHack.
var progressHackTests = []struct {
	status      webdm.Status
	scopeData   *ProgressHack
	expected    interface{}
	expectError bool
}{
	// Valid ProgressHacks
	{webdm.StatusUndefined, &ProgressHack{webdm.StatusUndefined}, &StorePreview{}, false},
	{webdm.StatusNotInstalled, &ProgressHack{webdm.StatusUndefined}, &StorePreview{}, false},
	{webdm.StatusInstalled, &ProgressHack{webdm.StatusUndefined}, &InstalledPreview{}, false},
	{webdm.StatusInstalling, &ProgressHack{webdm.StatusUndefined}, &StorePreview{}, false},
	{webdm.StatusUninstalling, &ProgressHack{webdm.StatusUndefined}, &StorePreview{}, false},

	{webdm.StatusUndefined, &ProgressHack{webdm.StatusInstalled}, &InstallingPreview{}, false},
	{webdm.StatusNotInstalled, &ProgressHack{webdm.StatusInstalled}, &InstallingPreview{}, false},
	{webdm.StatusInstalled, &ProgressHack{webdm.StatusInstalled}, &InstalledPreview{}, false},
	{webdm.StatusInstalling, &ProgressHack{webdm.StatusInstalled}, &InstallingPreview{}, false},
	{webdm.StatusUninstalling, &ProgressHack{webdm.StatusInstalled}, &InstallingPreview{}, false},

	{webdm.StatusUndefined, &ProgressHack{webdm.StatusNotInstalled}, &UninstallingPreview{}, false},
	{webdm.StatusNotInstalled, &ProgressHack{webdm.StatusNotInstalled}, &StorePreview{}, false},
	{webdm.StatusInstalled, &ProgressHack{webdm.StatusNotInstalled}, &UninstallingPreview{}, false},
	{webdm.StatusInstalling, &ProgressHack{webdm.StatusNotInstalled}, &UninstallingPreview{}, false},
	{webdm.StatusUninstalling, &ProgressHack{webdm.StatusNotInstalled}, &UninstallingPreview{}, false},

	// Invalid ProgressHacks
	{webdm.StatusUndefined, &ProgressHack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusUndefined, &ProgressHack{webdm.StatusUninstalling}, nil, true},

	{webdm.StatusNotInstalled, &ProgressHack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusNotInstalled, &ProgressHack{webdm.StatusUninstalling}, nil, true},

	{webdm.StatusInstalled, &ProgressHack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusInstalled, &ProgressHack{webdm.StatusUninstalling}, nil, true},

	{webdm.StatusInstalling, &ProgressHack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusInstalling, &ProgressHack{webdm.StatusUninstalling}, nil, true},

	{webdm.StatusUninstalling, &ProgressHack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusUninstalling, &ProgressHack{webdm.StatusUninstalling}, nil, true},
}

// Test typical NewPreview usage with temporary progress-hack related stuff.
func TestNewPreview_progressHack(t *testing.T) {
	for i, test := range progressHackTests {
		snap := webdm.Package{Status: test.status}
		metadata := scopes.NewActionMetadata("us", "phone")

		metadata.SetScopeData(test.scopeData)

		preview, err := NewPreview(snap, metadata)
		if test.expectError {
			if err == nil {
				t.Errorf("Test case %d: Expected an error", i)
			}
		} else {
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
}

// Test that ProgressHack cannot be retrieved if it's nil
func TestRetrieveProgressHack(t *testing.T) {
	metadata := scopes.NewActionMetadata("us", "phone")
	if retrieveProgressHack(metadata, nil) {
		t.Error("Expected retrieval to fail due to nil ProgressHack")
	}
}
