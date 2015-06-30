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
//
// Returns:
// - List of WebDM Package structs
// - Error (nil if none)
func GetPackageList(packageManager packages.Manager, department string) ([]webdm.Package, error) {
	var packages []webdm.Package
	var err error

	switch department {
	case "installed":
		packages, err = packageManager.GetInstalledPackages()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve installed packages: %s", err)
		}

	default:
		packages, err = packageManager.GetStorePackages()
		if err != nil {
			return nil, fmt.Errorf("Unable to retrieve store packages: %s", err)
		}
	}

	return packages, nil
}
