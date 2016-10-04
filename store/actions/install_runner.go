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
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// InstallRunner is an action Runner to handle the installation of a specific
// package.
type InstallRunner struct{}

// NewInstallRunner creates a new InstallRunner.
//
// Returns:
// - Pointer to new InstallRunner.
// - Error (nil if none).
func NewInstallRunner() (*InstallRunner, error) {
	return new(InstallRunner), nil
}

// Run installs the snap with the given ID.
//
// Parameters:
// packageManager: Package manager to use for installing the snap.
// snapId: ID of the snap to install.
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner InstallRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	objectPath, err := packageManager.Install(snapId)
	if err != nil {
		return nil, fmt.Errorf(`Unable to install package with ID "%s": %s`, snapId, err)
	}

	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	metadata := operation.Metadata{
		InstallRequested: true,
		ObjectPath:       objectPath,
	}

	response.SetScopeData(metadata)

	return response, nil
}
