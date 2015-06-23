package store

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// UninstallActionRunner is an ActionRunner to handle the uninstallation of a
// specific package.
type UninstallActionRunner struct{}

// NewUninstallActionRunner creates a new UninstallActionRunner.
//
// Returns:
// - Pointer to new UninstallActionRunner.
// - Error (nil if none).
func NewUninstallActionRunner() (*UninstallActionRunner, error) {
	return new(UninstallActionRunner), nil
}

// Run uninstalls the snap with the given ID.
//
// Parameters:
// packageManager: Package manager to use for uninstalling the snap.
// snapId: ID of the snap to uninstall.
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner UninstallActionRunner) Run(packageManager PackageManager, snapId string) (*scopes.ActivationResponse, error) {
	err := packageManager.Uninstall(snapId)
	if err != nil {
		return nil, fmt.Errorf(`Unable to uninstall package with ID "%s": %s`, snapId, err)
	}

	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	// Tell the preview when to stop showing the refresh page
	response.SetScopeData(ProgressHack{webdm.StatusNotInstalled})

	return response, nil
}
