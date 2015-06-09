package scope

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
)

// All possible actions in this scope
type ActionId int

const (
	ActionInstall ActionId = iota + 1
	ActionUninstall
	ActionOpen

	// Temporary actions for manual refresh
	ActionRefreshInstalling
	ActionRefreshUninstalling
	ActionOk
)

// ActionRunner is an interface for a factory of action handlers.
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
	case ActionRefreshInstalling:
		return NewRefreshInstallingActionRunner()
	case ActionRefreshUninstalling:
		return NewRefreshUninstallingActionRunner()
	case ActionOk:
		return NewOkActionRunner()
	default:
		var actionRunner ActionRunner
		return actionRunner, fmt.Errorf(`Unsupported action ID: "%d"`, actionId)
	}
}
