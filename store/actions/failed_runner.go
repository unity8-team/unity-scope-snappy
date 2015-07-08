package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// FailedRunner is an action Runner to handle a failed install or uninstall
// operation.
type FailedRunner struct{}

// NewFailedRunner creates a new FailedRunner.
//
// Returns:
// - Pointer to new FailedRunner.
// - Error (nil if none).
func NewFailedRunner() (*FailedRunner, error) {
	return new(FailedRunner), nil
}

// Run shoves the failed state into the metadata to be passed to the preview.
//
// Parameters:
// stateManager: Package state manager (not used).
// snapId: ID of the snap upon which the operation just had an error (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner FailedRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	response.SetScopeData(operation.Metadata{Failed: true})

	return response, nil
}
