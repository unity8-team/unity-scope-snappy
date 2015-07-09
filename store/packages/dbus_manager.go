package packages

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

// DbusManager is an interface to be implemented by any struct that supports
// asynchronous package installation/uninstallation via dbus in the manner used
// by this scope.
type DbusManager interface {
	Connect() error
	Install(packageId string) (dbus.ObjectPath, error)
	Uninstall(packageId string) (dbus.ObjectPath, error)
}
