package mocks

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

// MockBusObject is a mocked implementation of the dbus.BusObject interface, for
// use within tests.
type MockBusObject struct {
	CallCalled bool
	GoCalled bool
	GetPropertyCalled bool
	DestinationCalled bool
	PathCalled bool

	Method string
	Args []interface{}
}

func (mock *MockBusObject) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {
	mock.CallCalled = true
	mock.Method = method
	mock.Args = args
	return &dbus.Call{Body: []interface{}{dbus.ObjectPath("/foo/1")}}
}

func (mock *MockBusObject) Go(method string, flags dbus.Flags, ch chan *dbus.Call, args ...interface{}) *dbus.Call {
	mock.GoCalled = true
	mock.Method = method
	mock.Args = args
	return &dbus.Call{Body: []interface{}{dbus.ObjectPath("/foo/1")}}
}

func (mock *MockBusObject) GetProperty(p string) (dbus.Variant, error) {
	mock.GetPropertyCalled = true
	return dbus.MakeVariant("foo"), nil
}

func (mock *MockBusObject) Destination() string {
	mock.DestinationCalled = true
	return "foo"
}

func (mock *MockBusObject) Path() dbus.ObjectPath {
	mock.PathCalled = true
	return dbus.ObjectPath("/foo/1")
}
