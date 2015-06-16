package scope

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
)

type ActionId int

// All possible actions in this scope
const (
	ActionInstall ActionId = iota + 1
	ActionUninstall
	ActionOpen
)

// ActionRunner is an interface to be implemented by various action handlers
// throughout the scope.
type ActionRunner interface {
	Run(packageManager PackageManager, snapId string) (*scopes.ActivationResponse, error)
}

// NewActionRunner is a factory for getting the correct ActionRunner for a given
// ActionId.
//
// Parameters:
// actionId: The ID of the action needing to be handled.
func NewActionRunner(actionId ActionId) (ActionRunner, error) {
	switch actionId {
	case ActionInstall:
		return NewInstallActionRunner()
	case ActionUninstall:
		return NewUninstallActionRunner()
	case ActionOpen:
		return NewOpenActionRunner()
	default:
		return nil, fmt.Errorf(`Unsupported action ID: "%d"`, actionId)
	}
}
