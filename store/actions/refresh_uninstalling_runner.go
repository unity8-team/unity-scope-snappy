package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
	"launchpad.net/unity-scope-snappy/store/progress"
	"launchpad.net/unity-scope-snappy/webdm"
)

// RefreshUninstallingRunner is an action Runner to handle a manual refresh of
// the state of a specific package. This is a temporary workaround for the fact
// that actual progress isn't yet available for this scope.
type RefreshUninstallingRunner struct{}

// NewRefreshUninstallingRunner creates a new
// RefreshUninstallingRunner.
//
// Returns:
// - Pointer to new RefreshUninstallingRunner.
// - Error (nil if none).
func NewRefreshUninstallingRunner() (*RefreshUninstallingRunner, error) {
	return new(RefreshUninstallingRunner), nil
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
func (runner RefreshUninstallingRunner) Run(packageManager packages.Manager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	// Tell the preview when to stop showing the refresh page
	response.SetScopeData(progress.Hack{webdm.StatusNotInstalled})

	return response, nil
}
