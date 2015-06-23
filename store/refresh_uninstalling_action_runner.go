package store

import (
	"launchpad.net/unity-scope-snappy/store/packages"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// RefreshUninstallingActionRunner is an ActionRunner to handle a manual refresh
// of the state of a specific package. This is a temporary workaround for the
// fact that actual progress isn't yet available for this scope.
type RefreshUninstallingActionRunner struct{}

// NewRefreshUninstallingActionRunner creates a new
// RefreshUninstallingActionRunner.
//
// Returns:
// - Pointer to new RefreshUninstallingActionRunner.
// - Error (nil if none).
func NewRefreshUninstallingActionRunner() (*RefreshUninstallingActionRunner, error) {
	return new(RefreshUninstallingActionRunner), nil
}

// Run refreshes the prevew while passing along progress information.
//
// Parameters:
// packageManager: Package manager (not used).
// snapId: ID of the snap (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview
// - Error (nil if none)
func (runner RefreshUninstallingActionRunner) Run(packageManager packages.Manager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	// Tell the preview when to stop showing the refresh page
	response.SetScopeData(ProgressHack{webdm.StatusNotInstalled})

	return response, nil
}
