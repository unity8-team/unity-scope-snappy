package actions

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
)

type ActionId string

// All possible actions in this scope
const (
	ActionInstall          ActionId = "install"
	ActionUninstall                 = "uninstall"
	ActionUninstallConfirm          = "uninstall_confirm"
	ActionUninstallCancel           = "uninstall_cancel"
	ActionOpen                      = "open"

	// Temporary actions for manual refresh
	ActionRefreshInstalling   = "refresh_install"
	ActionRefreshUninstalling = "refresh_uninstall"
	ActionOk                  = "ok"

	// Actions from the progress widget
	ActionFinished = "finished"
	ActionFailed   = "failed"
)

// Runner is an interface to be implemented by the action handlers throughout
// the scope.
type Runner interface {
	Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error)
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
	case ActionUninstallConfirm:
		return NewConfirmUninstallRunner()
	case ActionUninstallCancel:
		return NewCancelUninstallRunner()
	case ActionOpen:
		return NewOpenRunner()
	case ActionRefreshInstalling:
		return NewRefreshInstallingRunner()
	case ActionRefreshUninstalling:
		return NewRefreshUninstallingRunner()
	case ActionOk:
		return NewOkRunner()

	// Actions from the progress widget
	case ActionFinished:
		return NewFinishedRunner()
	case ActionFailed:
		return NewFailedRunner()
	default:
		return nil, fmt.Errorf(`Unsupported action ID: "%s"`, actionId)
	}
}
