package store

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// OpenActionRunner is an ActionRunner to handle the launch of a
// specific snap.
type OpenActionRunner struct{}

// NewOpenActionRunner creates a new OpenActionRunner
//
// Returns:
// - Pointer to new OpenActionRunner
// - Error (nil if none)
func NewOpenActionRunner() (*OpenActionRunner, error) {
	return new(OpenActionRunner), nil
}

// Run is where a snap would be launched, if such a thing were supported.
//
// Parameters:
// packageManager: Package manager (not used).
// snapId: ID of the snap to launch (not used).
//
// Return:
// - A nil pointer to an ActivationResponse
// - An error saying that this isn't supported (yet).
func (runner OpenActionRunner) Run(packageManager packages.Manager, snapId string) (*scopes.ActivationResponse, error) {
	return nil, fmt.Errorf(`Unable to open package with ID "%s": Opening snaps is not yet supported`, snapId)
}
