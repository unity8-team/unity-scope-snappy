package daemon

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

type DbusServer struct {
	connection *dbus.Conn
}

func (server *DbusServer) Connect() error {
	var err error
	server.connection, err = dbus.SessionBus()
	return err
}

func (server *DbusServer) Names() []string {
	if server.connection == nil {
		return nil
	}

	return server.connection.Names()
}

func (server *DbusServer) RequestName(name string, flags dbus.RequestNameFlags) (dbus.RequestNameReply, error) {
	if server.connection == nil {
		return 0, fmt.Errorf("Server is not connected")
	}

	return server.connection.RequestName(name, flags)
}

func (server *DbusServer) GetNameOwner(name string) (string, error) {
	var owner string
	if server.connection == nil {
		return owner, fmt.Errorf("Server is not connected")
	}

	object := server.connection.BusObject()
	err := object.Call("org.freedesktop.DBus.GetNameOwner", 0, name).Store(&owner)
	return owner, err
}

func (server *DbusServer) Export(object interface{}, path dbus.ObjectPath, iface string) error {
	if server.connection == nil {
		return fmt.Errorf("Server is not connected")
	}

	return server.connection.Export(object, path, iface)
}

func (server *DbusServer) Emit(path dbus.ObjectPath, name string, values ...interface{}) error {
	if server.connection == nil {
		return fmt.Errorf("Server is not connected")
	}

	return server.connection.Emit(path, name, values...)
}
