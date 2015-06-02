package main

import (
	"launchpad.net/unity-scope-snappy/webdm"
)

// PackageManager is an interface to be implemented by any struct that supports
// the type of package management needed by this scope.
type PackageManager interface {
	GetInstalledPackages() ([]webdm.Package, error)
	GetStorePackages() ([]webdm.Package, error)
	Query(packageId string) (*webdm.Package, error)
	Install(packageId string) error
	Uninstall(packageId string) error
}
