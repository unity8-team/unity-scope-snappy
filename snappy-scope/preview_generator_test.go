package main

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
	{webdm.StatusUndefined, nil, "*main.StorePreview"},
	{webdm.StatusInstalled, nil, "*main.InstalledPreview"},
	{webdm.StatusNotInstalled, nil, "*main.StorePreview"},
	{webdm.StatusInstalling, nil, "*main.StorePreview"},
	{webdm.StatusUninstalling, nil, "*main.StorePreview"},

	// Temporary progress-related preview tests
	{webdm.StatusNotInstalled, &ProgressHack{webdm.StatusInstalled}, "*main.InstallingPreview"},
	{webdm.StatusInstalling, &ProgressHack{webdm.StatusInstalled}, "*main.InstallingPreview"},
	{webdm.StatusInstalled, &ProgressHack{webdm.StatusInstalled}, "*main.InstalledPreview"},
	{webdm.StatusInstalled, &ProgressHack{webdm.StatusNotInstalled}, "*main.UninstallingPreview"},
	{webdm.StatusUninstalling, &ProgressHack{webdm.StatusNotInstalled}, "*main.UninstallingPreview"},
	{webdm.StatusNotInstalled, &ProgressHack{webdm.StatusNotInstalled}, "*main.StorePreview"},
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
