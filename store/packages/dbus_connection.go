package packages

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

// DbusConnection is an interface of the DBus connection used by the DBus
// package manager client.
type DbusConnection interface {
	Names() []string
	Object(dest string, path dbus.ObjectPath) dbus.BusObject
}
