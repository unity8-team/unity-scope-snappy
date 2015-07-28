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

package utilities

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/store/packages"
	"launchpad.net/unity-scope-snappy/webdm"
)

// GetPackageList is used to obtain a package list for a specific department.
//
// Parameters:
// packageManager: Package manager to use to obtain package list.
// department: The department whose packages should be listed.
// query: Search query from the scope.
//
// Returns:
// - List of WebDM Package structs
// - Error (nil if none)
func GetPackageList(packageManager packages.WebdmManager, department string, query string) ([]webdm.Package, error) {
	var packages []webdm.Package
	var err error

	switch department {
	case "installed":
		packages, err = packageManager.GetInstalledPackages(query)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve installed packages: %s", err)
		}

	default:
		packages, err = packageManager.GetStorePackages(query)
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve store packages: %s", err)
		}
	}

	return packages, nil
}
