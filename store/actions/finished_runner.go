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

// FinishedRunner is an action Runner to handle the successful end of an install
// or uninstall operation.
type FinishedRunner struct{}

// NewFinishedRunner creates a new FinishedRunner.
//
// Returns:
// - Pointer to new FinishedRunner.
// - Error (nil if none).
func NewFinishedRunner() (*FinishedRunner, error) {
	return new(FinishedRunner), nil
}

// Run shoves the success state into the metadata to be passed to the preview.
//
// Parameters:
// stateManager: Package state manager (not used).
// snapId: ID of the snap upon which the operation just successfully completed
//         (not used).
//
// Return:
// - Pointer to an ActivationResponse for showing the preview.
// - Error (nil if none).
func (runner FinishedRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	response := scopes.NewActivationResponse(scopes.ActivationShowPreview)

	response.SetScopeData(operation.Metadata{Finished: true})

	return response, nil
}
