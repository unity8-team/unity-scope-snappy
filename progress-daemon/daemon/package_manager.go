package daemon

import (
	"launchpad.net/unity-scope-snappy/webdm"
)

// PackageManager is an interface to be implemented by any struct that supports
// the type of package management needed by this daemon.
type PackageManager interface {
	Query(packageId string) (*webdm.Package, error)
	Install(packageId string) error
	Uninstall(packageId string) error
}
