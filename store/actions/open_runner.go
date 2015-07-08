package actions

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// OpenRunner is an action Runner to handle the launch of a specific snap.
type OpenRunner struct{}

// NewOpenRunner creates a new OpenRunner
//
// Returns:
// - Pointer to new OpenRunner
// - Error (nil if none)
func NewOpenRunner() (*OpenRunner, error) {
	return new(OpenRunner), nil
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
func (runner OpenRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	return nil, fmt.Errorf(`Unable to open package with ID "%s": Opening snaps is not yet supported`, snapId)
}
