package daemon

import (
	"github.com/godbus/dbus"
)

// PackageManager is an interface to be implemented by any struct that supports
// the type of package management needed by this daemon.
type PackageManager interface {
	Install(packageId string) (dbus.ObjectPath, *dbus.Error)
	Uninstall(packageId string) (dbus.ObjectPath, *dbus.Error)
}
