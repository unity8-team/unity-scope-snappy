/* Copyright (C) 2015 Canonical Ltd.
 *
 * This file is part of unity-scope-snappy.
 *
 * unity-scope-snappy is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option) any
 * later version.
 *
 * unity-scope-snappy is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * unity-scope-snappy. If not, see <http://www.gnu.org/licenses/>.
 */

package actions

import (
	"fmt"
	"launchpad.net/go-unityscopes/v2"
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

	// Actions from the progress widget
	case ActionFinished:
		return NewFinishedRunner()
	case ActionFailed:
		return NewFailedRunner()
	default:
		return nil, fmt.Errorf(`Unsupported action ID: "%s"`, actionId)
	}
}
