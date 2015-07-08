package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// CancelUninstallRunner is an action Runner to handle the case when the
// uninstallation of a package is canceled.
type CancelUninstallRunner struct{}

// NewCancelUninstallRunner creates a new CancelUninstallRunner.
//
// Returns:
// - Pointer to new CancelUninstallRunner.
// - Error (nil if none).
func NewCancelUninstallRunner() (*CancelUninstallRunner, error) {
	return new(CancelUninstallRunner), nil
}

// Run simply refreshes the preview.
//
// Parameters:
// stateManager: Package state manager (not used).
// snapId: ID of the specific snap (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner CancelUninstallRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	return scopes.NewActivationResponse(scopes.ActivationShowPreview), nil
}
