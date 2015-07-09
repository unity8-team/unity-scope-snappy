package packages

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

const (
	defaultDbusObject          = "com.canonical.applications.WebdmPackageManager"
	defaultDbusObjectInterface = "com.canonical.applications.Download"
	defaultInstallMethod       = defaultDbusObjectInterface + ".Install"
	defaultUninstallMethod     = defaultDbusObjectInterface + ".Uninstall"
)

// DbusManagerClient is a DBus client for communicating with the WebDM Package
// Manager DBus service.
type DbusManagerClient struct {
	connection *dbus.Conn // Connection to the dbus bus

	dbusObject          string
	dbusObjectInterface string

	installMethod   string
	uninstallMethod string
}

// NewDbusManagerClient creates a new DbusManagerClient.
func NewDbusManagerClient() *DbusManagerClient {
	client := new(DbusManagerClient)

	client.dbusObject = defaultDbusObject
	client.dbusObjectInterface = defaultDbusObjectInterface

	client.installMethod = defaultInstallMethod
	client.uninstallMethod = defaultUninstallMethod

	return client
}

// Connect simply initializes a connection to the DBus session bus.
//
// Returns:
// - Error (nil if none)
func (client *DbusManagerClient) Connect() error {
	var err error
	client.connection, err = dbus.SessionBus()
	return err
}

// Install requests that the Package Manager service install the given package.
//
// Parameters:
// packageId: The ID of the package to install.
//
// Returns:
// - DBus object path to monitor the install operation.
// - Error (nil if none).
func (client *DbusManagerClient) Install(packageId string) (dbus.ObjectPath, error) {
	if client.connection == nil {
		return "", fmt.Errorf("Client is not connected")
	}

	busObject := client.connection.Object(client.dbusObject, "/")

	var objectPath dbus.ObjectPath
	err := busObject.Call(client.installMethod, 0, packageId).Store(&objectPath)

	return objectPath, err
}

// Uninstall requests that the Package Manager service uninstall the given
// package.
//
// Parameters:
// packageId: The ID of the package to uninstall.
//
// Returns:
// - DBus object path to monitor the uninstall operation.
// - Error (nil if none).
func (client *DbusManagerClient) Uninstall(packageId string) (dbus.ObjectPath, error) {
	if client.connection == nil {
		return "", fmt.Errorf("Client is not connected")
	}

	busObject := client.connection.Object(client.dbusObject, "/")

	var objectPath dbus.ObjectPath
	err := busObject.Call(client.uninstallMethod, 0, packageId).Store(&objectPath)

	return objectPath, err
}
