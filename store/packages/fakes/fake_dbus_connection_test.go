package fakes

import (
	"launchpad.net/unity-scope-snappy/store/packages/mocks"
	"reflect"
	"testing"
)

// Test that Names returns a single name
func TestFakeDbusConnection_names(t *testing.T) {
	connection := FakeDbusConnection{}

	names := connection.Names()
	if len(names) != 1 {
		t.Fatalf("Got %d names, expected 1", len(names))
	}
}

// Test that Object simply returns the given DbusObject.
func TestFakeDbusConnection_object(t *testing.T) {
	mock := &mocks.MockBusObject{}
	connection := FakeDbusConnection{mock}

	if !reflect.DeepEqual(connection.Object("foo", "bar"), mock) {
		t.Error("Expected the fake connection to return the given mock")
	}
}
