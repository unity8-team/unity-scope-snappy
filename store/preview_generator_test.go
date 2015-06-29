package store

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/progress"
	"launchpad.net/unity-scope-snappy/webdm"
	"reflect"
	"testing"
)

// Data for TestNewPreview.
var newPreviewTests = []struct {
	status    webdm.Status
	scopeData *progress.Hack
	expected  interface{}
}{
	{webdm.StatusUndefined, nil, &PackagePreview{}},
	{webdm.StatusInstalled, nil, &PackagePreview{}},
	{webdm.StatusNotInstalled, nil, &PackagePreview{}},
	{webdm.StatusInstalling, nil, &PackagePreview{}},
	{webdm.StatusUninstalling, nil, &PackagePreview{}},
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
	scopeData   *progress.Hack
	expected    interface{}
	expectError bool
}{
	// Valid ProgressHacks
	{webdm.StatusUndefined, &progress.Hack{webdm.StatusUndefined}, &PackagePreview{}, false},
	{webdm.StatusNotInstalled, &progress.Hack{webdm.StatusUndefined}, &PackagePreview{}, false},
	{webdm.StatusInstalled, &progress.Hack{webdm.StatusUndefined}, &PackagePreview{}, false},
	{webdm.StatusInstalling, &progress.Hack{webdm.StatusUndefined}, &PackagePreview{}, false},
	{webdm.StatusUninstalling, &progress.Hack{webdm.StatusUndefined}, &PackagePreview{}, false},

	{webdm.StatusUndefined, &progress.Hack{webdm.StatusInstalled}, &InstallingPreview{}, false},
	{webdm.StatusNotInstalled, &progress.Hack{webdm.StatusInstalled}, &InstallingPreview{}, false},
	{webdm.StatusInstalled, &progress.Hack{webdm.StatusInstalled}, &PackagePreview{}, false},
	{webdm.StatusInstalling, &progress.Hack{webdm.StatusInstalled}, &InstallingPreview{}, false},
	{webdm.StatusUninstalling, &progress.Hack{webdm.StatusInstalled}, &InstallingPreview{}, false},

	{webdm.StatusUndefined, &progress.Hack{webdm.StatusNotInstalled}, &UninstallingPreview{}, false},
	{webdm.StatusNotInstalled, &progress.Hack{webdm.StatusNotInstalled}, &PackagePreview{}, false},
	{webdm.StatusInstalled, &progress.Hack{webdm.StatusNotInstalled}, &UninstallingPreview{}, false},
	{webdm.StatusInstalling, &progress.Hack{webdm.StatusNotInstalled}, &UninstallingPreview{}, false},
	{webdm.StatusUninstalling, &progress.Hack{webdm.StatusNotInstalled}, &UninstallingPreview{}, false},

	// Invalid combinations
	{webdm.StatusUndefined, &progress.Hack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusUndefined, &progress.Hack{webdm.StatusUninstalling}, nil, true},

	{webdm.StatusNotInstalled, &progress.Hack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusNotInstalled, &progress.Hack{webdm.StatusUninstalling}, nil, true},

	{webdm.StatusInstalled, &progress.Hack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusInstalled, &progress.Hack{webdm.StatusUninstalling}, nil, true},

	{webdm.StatusInstalling, &progress.Hack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusInstalling, &progress.Hack{webdm.StatusUninstalling}, nil, true},

	{webdm.StatusUninstalling, &progress.Hack{webdm.StatusInstalling}, nil, true},
	{webdm.StatusUninstalling, &progress.Hack{webdm.StatusUninstalling}, nil, true},
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
