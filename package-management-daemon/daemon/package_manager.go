package daemon

import (
	"github.com/godbus/dbus"
	"github.com/snapcore/snapd/client"
)

// PackageManager is an interface to be implemented by any struct that supports
// the type of package management needed by this daemon.
type PackageManager interface {
	Query(packageId string) (*client.Snap, error)
	Install(packageId string) (dbus.ObjectPath, *dbus.Error)
	Uninstall(packageId string) (dbus.ObjectPath, *dbus.Error)
}
