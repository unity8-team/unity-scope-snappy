package actions

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
)

type ActionId int

// All possible actions in this scope
const (
	ActionInstall ActionId = iota + 1
	ActionUninstall
	ActionOpen

	// Temporary actions for manual refresh
	ActionRefreshInstalling
	ActionRefreshUninstalling
	ActionOk
)

// Runner is an interface to be implemented by the action handlers throughout
// the scope.
type Runner interface {
	Run(packageManager packages.Manager, snapId string) (*scopes.ActivationResponse, error)
}

// NewRunner is a factory for getting the correct Runner for a given ActionId.
//
// Parameters:
// actionId: The ID of the action needing to be handled.
func NewRunner(actionId ActionId) (Runner, error) {
	switch actionId {
	case ActionInstall:
		return NewInstallRunner()
	case ActionUninstall:
		return NewUninstallRunner()
	case ActionOpen:
		return NewOpenRunner()
	case ActionRefreshInstalling:
		return NewRefreshInstallingRunner()
	case ActionRefreshUninstalling:
		return NewRefreshUninstallingRunner()
	case ActionOk:
		return NewOkRunner()
	default:
		return nil, fmt.Errorf(`Unsupported action ID: "%d"`, actionId)
	}
}
