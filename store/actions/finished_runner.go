package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// FinishedRunner is an action Runner to handle the successful end of an install
// or uninstall operation.
type FinishedRunner struct{}

// NewFinishedRunner creates a new FinishedRunner.
//
// Returns:
// - Pointer to new FinishedRunner.
// - Error (nil if none).
func NewFinishedRunner() (*FinishedRunner, error) {
	return new(FinishedRunner), nil
}

// Run shoves the success state into the metadata to be passed to the preview.
//
// Parameters:
// stateManager: Package state manager (not used).
// snapId: ID of the snap upon which the operation just successfully completed
//         (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner FinishedRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	response.SetScopeData(operation.Metadata{Finished: true})

	return response, nil
}
