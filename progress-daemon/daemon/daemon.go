package daemon

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus/introspect"
)

const (
	// busName is the name to be requested from the DBus session bus.
	busName = "com.canonical.applications.WebdmPackageManager"

	// introspectionXml is the XML to be used for the Introspection interface.
	introspectionXml = `
		<node>
			<interface name="com.canonical.applications.WebdmPackageManager">
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
	packageManager *WebdmPackageManagerInterface
}

// New creates a new Daemon setup to poll WebDM at a specific URL.
//
// Parameters:
// webdmApiUrl: WebDM API URL to poll.
//
// Returns:
// - New daemon
// - Error (nil if none)
func New(webdmApiUrl string) (*Daemon, error) {
	daemon := new(Daemon)

	daemon.server = new(DbusServer)

	var err error
	daemon.packageManager, err = NewWebdmPackageManagerInterface(daemon.server, webdmApiUrl)
	if err != nil {
		return nil, fmt.Errorf(`Unable to create package manager interface with API URL "%s"`, webdmApiUrl)
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

	err = daemon.server.Export(daemon.packageManager, "/", "com.canonical.applications.WebdmPackageManager")
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
