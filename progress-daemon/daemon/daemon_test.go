package daemon

import (
	"testing"
)

// Test typical New usage.
func TestNew(t *testing.T) {
	daemon, err := New("")
	if err != nil {
		t.Fatalf("Unexpected error when creating daemon: %s", err)
	}

	if daemon.server == nil {
		t.Error("Daemon dbus server was unexpectedly nil")
	}

	if daemon.packageManager == nil {
		t.Errorf("Daemon package manager was unexpectedly nil")
	}
}

// Test that New fails if given an invalid API URL
func TestNew_invalidUrl(t *testing.T) {
	_, err := New(":")
	if err == nil {
		t.Error("Expected an error due to invalid API URL")
	}
}

// Test typical Run usage.
func TestDaemonRunStop(t *testing.T) {
	daemon, err := New("")
	if err != nil {
		t.Fatalf("Unexpected error when creating daemon: %s", err)
	}

	err = daemon.Run()
	if err != nil {
		t.Error("Unexpected error while running daemon")
	}

	// Obtain our own unique name
	names := daemon.server.Names()
	if len(names) < 1 {
		t.Errorf("Got %d names, expected at least 1", len(names))
	}

	// Make sure we're running by making sure we own the name we expect
	name, err := daemon.server.GetNameOwner(busName)
	if err != nil {
		t.Error("Unexpected error while requesting name owner")
	}

	// Make sure we own that name
	if name != names[0] {
		t.Errorf(`Name owner was "%s", expected the owner to be us ("%s")`, name, names[0])
	}
}

// Test dbus connection failure
func TestRun_connectionFailure(t *testing.T) {
	daemon, err := New("")
	if err != nil {
		t.Fatalf("Unexpected error when creating daemon: %s", err)
	}

	// Use fake server instead of real one
	daemon.server = &FakeDbusServer{failConnect: true}

	err = daemon.Run()
	if err == nil {
		t.Error("Expected an error due to failure to connect")
	}
}

// Test dbus name request failure
func TestRun_nameRequestFailure(t *testing.T) {
	daemon, err := New("")
	if err != nil {
		t.Fatalf("Unexpected error when creating daemon: %s", err)
	}

	// Use fake server instead of real one
	daemon.server = &FakeDbusServer{failRequestName: true}

	err = daemon.Run()
	if err == nil {
		t.Error("Expected an error due to failure to request name")
	}
}

// Test dbus name already taken
func TestRun_nameTaken(t *testing.T) {
	daemon, err := New("")
	if err != nil {
		t.Fatalf("Unexpected error when creating daemon: %s", err)
	}

	// Use fake server instead of real one
	daemon.server = &FakeDbusServer{nameAlreadyTaken: true}

	err = daemon.Run()
	if err == nil {
		t.Error("Expected an error due to name already being taken")
	}
}

// Test dbus introspection export failure
func TestRun_introspectionExportFailure(t *testing.T) {
	daemon, err := New("")
	if err != nil {
		t.Fatalf("Unexpected error when creating daemon: %s", err)
	}

	// Use fake server instead of real one
	daemon.server = &FakeDbusServer{
		failExport:                  true,
		failSpecificExportInterface: "org.freedesktop.DBus.Introspectable",
	}

	err = daemon.Run()
	if err == nil {
		t.Error("Expected an error due to failure to export")
	}
}

// Test dbus package manager export failure
func TestRun_packageManagerExportFailure(t *testing.T) {
	daemon, err := New("")
	if err != nil {
		t.Fatalf("Unexpected error when creating daemon: %s", err)
	}

	// Use fake server instead of real one
	daemon.server = &FakeDbusServer{
		failExport:                  true,
		failSpecificExportInterface: "com.canonical.applications.WebdmPackageManager",
	}

	err = daemon.Run()
	if err == nil {
		t.Error("Expected an error due to failure to export")
	}
}
