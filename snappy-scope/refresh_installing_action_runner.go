package main

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/webdm"
)

// RefreshInstallingActionRunner is an ActionRunner to handle a manual refresh
// of the state of a specific package. This is a temporary workaround for the
// fact that actual progress isn't yet available for this scope.
type RefreshInstallingActionRunner struct{}

// NewRefreshInstallingActionRunner creates a new RefreshInstallingActionRunner.
//
// Returns:
// - Pointer to new RefreshInstallingActionRunner.
// - Error (nil if none).
func NewRefreshInstallingActionRunner() (*RefreshInstallingActionRunner, error) {
	return new(RefreshInstallingActionRunner), nil
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
func (runner RefreshInstallingActionRunner) Run(packageManager PackageManager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	// Tell the preview when to stop showing the refresh page
	response.SetScopeData(ProgressHack{webdm.StatusInstalled})

	return response, nil
}
