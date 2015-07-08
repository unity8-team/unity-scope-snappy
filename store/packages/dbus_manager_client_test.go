package packages

import (
	"testing"
)

// Test typical NewDbusManagerClient usage.
func TestNewDbusManagerClient(t *testing.T) {
	client := NewDbusManagerClient()
	if client == nil {
		t.Fatal("Client was unexpectedly nil")
	}

	if client.dbusObject == "" {
		t.Error("Client didn't configure a dbus object")
	}

	if client.dbusObjectInterface == "" {
		t.Error("Client didn't configure a dbus object interface")
	}
}

// Test typical Connect usage.
func TestDbusManagerClient_connect(t *testing.T) {
	client := NewDbusManagerClient()
	err := client.Connect()
	if err != nil {
		t.Errorf("Unexpected error while connecting: %s", err)
	}

	names := client.connection.Names()
	if len(names) < 1 {
		t.Errorf("Got %d names, expected at least 1", len(names))
	}
}

// Test that trying to install before connecting results in an error.
func TestDbusManagerClient_install_beforeConnect(t *testing.T) {
	client := NewDbusManagerClient()
	_, err := client.Install("foo")
	if err == nil {
		t.Error("Expected an error due to install before connect")
	}
}

// Test that trying to uninstall before connecting results in an error.
func TestDbusManagerClient_uninstall_beforeConnect(t *testing.T) {
	client := NewDbusManagerClient()
	_, err := client.Uninstall("foo")
	if err == nil {
		t.Error("Expected an error due to uninstall before connect")
	}
}
