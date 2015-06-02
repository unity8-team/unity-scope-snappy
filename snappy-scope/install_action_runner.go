package main

import (
	"fmt"
	"launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// InstallActionRunner is an ActionRunner to handle the installation of a
// specific package.
type InstallActionRunner struct {}

// NewInstallActionRunner creates a new InstallActionRunner.
//
// Returns:
// - Pointer to new InstallActionRunner.
// - Error (nil if none).
func NewInstallActionRunner() (*InstallActionRunner, error) {
	return new(InstallActionRunner), nil
}

// Run installs the snap with the given ID.
//
// Parameters:
// packageManager: Package manager to use for installing the snap.
// snapId: ID of the snap to install.
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner InstallActionRunner) Run(packageManager PackageManager, snapId string) (*scopes.ActivationResponse, error) {
	err := packageManager.Install(snapId)
	if err != nil {
		return nil, fmt.Errorf(`Unable to install package with ID "%s": %s`, snapId, err)
	}

	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	// Tell the preview when to stop showing the refresh page
	response.SetScopeData(ProgressHack{webdm.StatusInstalled})

	return response, nil
}
