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
	"launchpad.net/unity-scope-snappy/store/packages"
)

// CancelUninstallRunner is an action Runner to handle the case when the
// uninstallation of a package is canceled.
type CancelUninstallRunner struct{}

// NewCancelUninstallRunner creates a new CancelUninstallRunner.
//
// Returns:
// - Pointer to new CancelUninstallRunner.
// - Error (nil if none).
func NewCancelUninstallRunner() (*CancelUninstallRunner, error) {
	return new(CancelUninstallRunner), nil
}

// Run simply refreshes the preview.
//
// Parameters:
// stateManager: Package state manager (not used).
// snapId: ID of the specific snap (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner CancelUninstallRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	return scopes.NewActivationResponse(scopes.ActivationShowPreview), nil
}
