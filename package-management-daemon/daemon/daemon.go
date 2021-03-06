package daemon

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

const (
	// busName is the name to be requested from the DBus session bus.
	busName = "com.canonical.applications.WebdmPackageManager"

	// baseObjectPath is the base object path to use for operation
	// notifications (dbus signals)
	baseObjectPath = "/com/canonical/applications/WebdmPackageManager/operation"

	// interfaceName is the name of the interface being implemented here.
	// The driver for this name is Unity8's QML Progress Widget, which is
	// hard-coded to look for this interface.
	interfaceName = "com.canonical.applications.Download"

	// introspectionXml is the XML to be used for the Introspection interface.
	introspectionXml = `
		<node>
			<interface name="` + interfaceName + `">
				<method name="Install">
					<arg name="packageId" type="s" direction="in"/>
				</method>
				<method name="Uninstall">
					<arg name="packageId" type="s" direction="in"/>
				</method>
				<signal name="progress">
					<arg name="received" type="t" />
					<arg name="total" type="t" />
				</signal>
				<signal name="finished">
					<arg name="path" type="s" />
				</signal>
				<signal name="error">
					<arg name="error" type="s" />
				</signal>
			</interface>` +
		introspect.IntrospectDataString +
		`</node>`
)

// Daemon represents the actual progress daemon.
type Daemon struct {
	server         DbusWrapper
	packageManager PackageManager
}

// New creates a new Daemon setup to poll WebDM at a specific URL.
//
// Parameters:
// webdmApiUrl: WebDM API URL to poll.
//
// Returns:
// - New daemon
// - Error (nil if none)
func New() (*Daemon, error) {
	daemon := new(Daemon)

	daemon.server = new(DbusServer)

	// dbusConnection: Connection to the dbus bus.
	// interfaceName: DBus interface name to implement.
	// baseObjectPath: Base object path to use for signals.
	// apiUrl: WebDM API URL.

	var err error
	daemon.packageManager, err = NewSnapdPackageManagerInterface(daemon.server,
		interfaceName, baseObjectPath)
	if err != nil {
		return nil, fmt.Errorf(`Unable to create package manager interface:"`, err)
	}

	return daemon, nil
}

// Run connects to the DBus session bus and prepares for receiving requests.
//
// Returns:
// - Error (nil if none)
func (daemon *Daemon) Run() error {
	err := daemon.server.Connect()
	if err != nil {
		return fmt.Errorf("Unable to connect: %s", err)
	}

	err = daemon.server.Export(introspect.Introspectable(introspectionXml), "/",
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return fmt.Errorf("Unable to export introspection: %s", err)
	}

	err = daemon.server.Export(daemon.packageManager, "/", interfaceName)
	if err != nil {
		return fmt.Errorf("Unable to export package manager interface: %s", err)
	}

	// Now that all interfaces are exported and ready, request our name. Things
	// are done in this order so that our interfaces aren't called before
	// they're exported.
	reply, err := daemon.server.RequestName(busName, dbus.NameFlagDoNotQueue)
	if err != nil {
		return fmt.Errorf(`Unable to get requested name "%s": %s`, busName, err)
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return fmt.Errorf(`Requested name "%s" was already taken`, busName)
	}

	return nil
}
