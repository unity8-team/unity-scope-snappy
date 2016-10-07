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

package webdm

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
func (snapd *SnapdClient) GetStorePackages(query string) ([]Package, error) {
	if query == "" {
		query = "."
	}
	snaps, _, err := snapd.snapdClient.Find(&client.FindOptions{
		Query: query,
	})
	if err != nil {
		return nil, fmt.Errorf("snapd: Error getting store packages: %s", err)
	}

	packages := make([]Package, 0)
	for _, snap := range snaps {
		snappkg := &Package{
			Id:            snap.ID,
			Name:          snap.Name,
			Version:       snap.Version,
			Type:          snap.Type,
			IconUrl:       snap.Icon,
			Description:   snap.Description,
			DownloadSize:  snap.DownloadSize,
			InstalledSize: snap.InstalledSize,
			Vendor:        snap.Developer,
		}
		packages = append(packages, *snappkg)
	}
	return packages, nil
}

func (snapd *SnapdClient) Query(snapName string) (*Package, error) {
	// Check first if the snap in question is already installed
	pkgs, err := snapd.snapdClient.List([]string{snapName})
	if err != nil {
		return nil, fmt.Errorf("snapd: Error getting package: %s", err)
	}

	var pkg *client.Snap
	if pkgs == nil {
		// The snap we're looking for is not installed, find it in the store
		pkg, _, err = snapd.snapdClient.FindOne(snapName)
		if err != nil {
			return nil, fmt.Errorf("snapd: Error getting package: %s", err)
		}
	} else {
		// The snap we're looking for is installed, use the local snap's info
		pkg = pkgs[0]
	}

	snap := new(Package)
	snap.Id = pkg.ID
	snap.Name = pkg.Name
	snap.Version = pkg.Version
	snap.Type = pkg.Type
	snap.IconUrl = pkg.Icon
	snap.Description = pkg.Description
	snap.DownloadSize = pkg.DownloadSize
	snap.InstalledSize = pkg.InstalledSize
	snap.Vendor = pkg.Developer

	if pkg.Status == client.StatusAvailable {
		snap.Status = StatusNotInstalled
	} else if pkg.Status == client.StatusInstalled {
		snap.Status = StatusInstalled
	} else if pkg.Status == client.StatusActive {
		snap.Status = StatusInstalled
	} else if pkg.Status == client.StatusRemoved {
		snap.Status = StatusNotInstalled
	}

	return snap, nil
}

func (snapd *SnapdClient) Install(packageId string) error {
	return nil
}

func (snapd *SnapdClient) Uninstall(packageId string) error {
	return nil
}
