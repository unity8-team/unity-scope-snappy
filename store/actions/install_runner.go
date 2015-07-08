package actions

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
	"launchpad.net/unity-scope-snappy/store/progress"
	"launchpad.net/unity-scope-snappy/webdm"
)

// InstallRunner is an action Runner to handle the installation of a specific
// package.
type InstallRunner struct{}

// NewInstallRunner creates a new InstallRunner.
//
// Returns:
// - Pointer to new InstallRunner.
// - Error (nil if none).
func NewInstallRunner() (*InstallRunner, error) {
	return new(InstallRunner), nil
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
func (runner InstallRunner) Run(packageManager packages.WebdmManager, snapId string) (*scopes.ActivationResponse, error) {
	err := packageManager.Install(snapId)
	if err != nil {
		return nil, fmt.Errorf(`Unable to install package with ID "%s": %s`, snapId, err)
	}

	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	// Tell the preview when to stop showing the refresh page
	response.SetScopeData(progress.Hack{webdm.StatusInstalled})

	return response, nil
}
