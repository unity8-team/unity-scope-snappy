/* Copyright (C) 2015-2016 Canonical Ltd.
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

package packages

import (
	"github.com/snapcore/snapd/client"
)

// WebdmManager is an interface to be implemented by any struct that supports
// the type of package management needed by this scope.
type WebdmManager interface {
	GetInstalledPackages() (map[string]struct{})
	GetStorePackages(query string) ([]client.Snap, error)
	Query(packageId string) (*client.Snap, error)
	Install(packageId string) error
	Uninstall(packageId string) error
}
