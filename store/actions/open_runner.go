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
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/packages"
)

// OpenRunner is an action Runner to handle the launch of a specific snap.
type OpenRunner struct{}

// NewOpenRunner creates a new OpenRunner
//
// Returns:
// - Pointer to new OpenRunner
// - Error (nil if none)
func NewOpenRunner() (*OpenRunner, error) {
	return new(OpenRunner), nil
}

// Run is where a snap would be launched, if such a thing were supported.
//
// Parameters:
// packageManager: Package manager (not used).
// snapId: ID of the snap to launch (not used).
//
// Return:
// - A nil pointer to an ActivationResponse
// - An error saying that this isn't supported (yet).
func (runner OpenRunner) Run(packageManager packages.DbusManager, snapId string) (*scopes.ActivationResponse, error) {
	return nil, fmt.Errorf(`Unable to open package with ID "%s": Opening snaps is not yet supported`, snapId)
}
