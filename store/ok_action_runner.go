package store

import (
	"launchpad.net/unity-scope-snappy/store/packages"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
)

// OkActionRunner is an ActionRunner to handle the "Ok" button.
type OkActionRunner struct{}

// NewOkActionRunner creates a new OkActionRunner
//
// Returns:
// - Pointer to new OkActionRunner
// - Error (nil if none)
func NewOkActionRunner() (*OkActionRunner, error) {
	return new(OkActionRunner), nil
}

// Run simply returns an ActivationResponse to show the preview.
//
// Parameters:
// packageManager: Package manager (not used).
// snapId: ID of the snap (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview
// - Error (nil if none)
func (runner OkActionRunner) Run(packageManager packages.Manager, snapId string) (*scopes.ActivationResponse, error) {
	return scopes.NewActivationResponse(scopes.ActivationShowPreview), nil
}
