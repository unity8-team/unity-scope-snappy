package actions

import (
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// OkRunner is an action Runner to handle the "Ok" button.
type OkRunner struct{}

// NewOkRunner creates a new OkRunner
//
// Returns:
// - Pointer to new OkRunner
// - Error (nil if none)
func NewOkRunner() (*OkRunner, error) {
	return new(OkRunner), nil
}

// Run simply returns an ActivationResponse to show the preview.
//
// Parameters:
// packageManager: Package manager (not used).
// snapId: ID of the snap (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner OkRunner) Run(packageManager packages.Manager, snapId string) (*scopes.ActivationResponse, error) {
	return scopes.NewActivationResponse(scopes.ActivationShowPreview), nil
}
