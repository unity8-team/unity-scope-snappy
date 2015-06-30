package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// ErrorRunner is an action Runner to handle an error during an install or
// uninstall operation.
type ErrorRunner struct{}

// NewErrorRunner creates a new ErrorRunner.
//
// Returns:
// - Pointer to new ErrorRunner.
// - Error (nil if none).
func NewErrorRunner() (*ErrorRunner, error) {
	return new(ErrorRunner), nil
}

// Run shoves the error state into the metadata to be passed to the preview.
//
// Parameters:
// stateManager: Package state manager (not used).
// snapId: ID of the snap upon which the operation just had an error (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner ErrorRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	response.SetScopeData(operation.Metadata{Error: true})

	return response, nil
}
