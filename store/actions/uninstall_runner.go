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
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// UninstallRunner is an action Runner to handle the uninstallation of a
// specific package.
type UninstallRunner struct{}

// NewUninstallRunner creates a new UninstallRunner.
//
// Returns:
// - Pointer to new UninstallRunner.
// - Error (nil if none).
func NewUninstallRunner() (*UninstallRunner, error) {
	return new(UninstallRunner), nil
}

// Run uninstalls the snap with the given ID.
//
// Parameters:
// packageManager: Package manager to use for uninstalling the snap.
// snapId: ID of the snap to uninstall.
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner UninstallRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	metadata := operation.Metadata{
		UninstallRequested: true,
	}

	response.SetScopeData(metadata)

	return response, nil
}
