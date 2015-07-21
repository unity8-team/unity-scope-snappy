package packages

import (
	"launchpad.net/unity-scope-snappy/webdm"
)

// WebdmManager is an interface to be implemented by any struct that supports
// the type of package management needed by this scope.
type WebdmManager interface {
	GetInstalledPackages(query string) ([]webdm.Package, error)
	GetStorePackages(query string) ([]webdm.Package, error)
	Query(packageId string) (*webdm.Package, error)
	Install(packageId string) error
	Uninstall(packageId string) error
}
