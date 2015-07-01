package operation

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

// Metadata exists to persist state throughout scope calls, even if the scope is
// killed.
type Metadata struct {
	InstallRequested   bool
	UninstallRequested bool
	UninstallConfirmed bool

	Finished bool
	Failed   bool

	ObjectPath dbus.ObjectPath
}
