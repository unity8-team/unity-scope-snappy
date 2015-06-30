package actions

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
	"launchpad.net/unity-scope-snappy/store/progress"
	"launchpad.net/unity-scope-snappy/webdm"
)

// UninstallRunner is an action Runner to handle the uninstallation of a
// specific package.
type UninstallRunner struct{}

// NewUninstallRunner creates a new UninstallRunner.
//
// Returns:
// - Pointer to new UninstallRunner.
// - Error (nil if none).
func NewUninstallRunner() (*UninstallRunner, error) {
	return new(UninstallRunner), nil
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
func (runner UninstallRunner) Run(packageManager packages.WebdmManager, snapId string) (*scopes.ActivationResponse, error) {
	err := packageManager.Uninstall(snapId)
	if err != nil {
		return nil, fmt.Errorf(`Unable to uninstall package with ID "%s": %s`, snapId, err)
	}

	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	// Tell the preview when to stop showing the refresh page
	response.SetScopeData(progress.Hack{webdm.StatusNotInstalled})

	return response, nil
}
