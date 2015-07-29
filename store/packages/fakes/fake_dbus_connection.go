package fakes

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

// DbusConnection is an interface of the DBus connection used by the DBus
// package manager client.
type FakeDbusConnection struct {
	DbusObject dbus.BusObject
}

func (fake FakeDbusConnection) Names() []string {
	return []string{":1.42"}
}

func (fake FakeDbusConnection) Object(dest string, path dbus.ObjectPath) dbus.BusObject {
	return fake.DbusObject
}
