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

// FailedRunner is an action Runner to handle a failed install or uninstall
// operation.
type FailedRunner struct{}

// NewFailedRunner creates a new FailedRunner.
//
// Returns:
// - Pointer to new FailedRunner.
// - Error (nil if none).
func NewFailedRunner() (*FailedRunner, error) {
	return new(FailedRunner), nil
}

// Run shoves the failed state into the metadata to be passed to the preview.
//
// Parameters:
// stateManager: Package state manager (not used).
// snapId: ID of the snap upon which the operation just had an error (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner FailedRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	response.SetScopeData(operation.Metadata{Failed: true})

	return response, nil
}
