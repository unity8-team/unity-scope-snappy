package daemon

import "launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"

// DbusWrapper is an interface to be satisfied by any struct that wants to be
// injectable into this daemon for dbus communication.
type DbusWrapper interface {
	Connect() error
	Names() []string
	RequestName(name string, flags dbus.RequestNameFlags) (dbus.RequestNameReply, error)
	GetNameOwner(name string) (string, error)
	Export(object interface{}, path dbus.ObjectPath, iface string) error
	Emit(path dbus.ObjectPath, name string, values ...interface{}) error
}
