/* Copyright (C) 2016 Canonical Ltd.
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
	"fmt"

	"github.com/snapcore/snapd/client"
)

// Client is the main struct allowing for communication with the webdm API.
type SnapdClient struct {
	snapdClientConfig client.Config
	snapdClient       *client.Client
}

// NewClient creates a new client for communicating with the webdm API
//
// Parameters:
// apiUrl: URL where WebDM is listening (host[:port])
//
// Returns:
// - Pointer to new client
// - Error (nil if none)
func NewSnapdClient() (*SnapdClient, error) {
	snapd := &SnapdClient{}
	snapd.snapdClient = client.New(&snapd.snapdClientConfig)

	return snapd, nil
}

// GetInstalledPackages sends an API request for a list of installed packages.
//
// Parameters:
// query: Search query for list.
//
// Returns:
// - Slice of Packags structs
// - Error (nil of none)
func (snapd *SnapdClient) GetInstalledPackages() (map[string]struct{}) {
	snaps, err := snapd.snapdClient.List(nil)
	if err != nil {
		fmt.Printf("snapd: Error getting installed packages: %s", err)
	}

	packages := make(map[string]struct{}, 0)
	for _, snap := range snaps {
		packages[snap.Name] = struct{}{}
	}
	return packages
}

// GetStorePackages sends an API request for a list of all packages in the
// store (including installed packages).
//
// Parameters:
// query: Search query for list.
//
// Returns:
// - Slice of Packags structs
// - Error (nil of none)
func (snapd *SnapdClient) GetStorePackages(query string) ([]client.Snap, error) {
	if query == "" {
		query = "."
	}
	snaps, _, err := snapd.snapdClient.Find(&client.FindOptions{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("snapd: Error getting store packages: %s", err)
	}

	packages := make([]client.Snap, 0)
	for _, snap := range snaps {
		if snap.Type != client.TypeApp {
			continue
		}
		packages = append(packages, *snap)
	}
	return packages, nil
}

func (snapd *SnapdClient) Query(snapName string) (*client.Snap, error) {
	// Check first if the snap in question is already installed
	pkg, _, err := snapd.snapdClient.Snap(snapName)
	if err != nil {
		pkg, _, err = snapd.snapdClient.FindOne(snapName)
		if err != nil {
			return nil, fmt.Errorf("snapd: Error getting package: %s", err)
		}
	}

	return pkg, nil
}

func (snapd *SnapdClient) Install(packageId string) error {
	return nil
}

func (snapd *SnapdClient) Uninstall(packageId string) error {
	return nil
}
