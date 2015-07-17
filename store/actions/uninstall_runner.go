package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages"
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
func (runner UninstallRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	metadata := operation.Metadata{
		UninstallRequested: true,
	}

	response.SetScopeData(metadata)

	return response, nil
}
