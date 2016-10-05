/* Copyright (C) 2016 Canonical Ltd.
 *
 * This file is part of unity-scope-snappy.
 *
 * unity-scope-snappy is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by the Free
 * Software Foundation, either version 3 of the License, or (at your option) any
 * later version.
 *
 * unity-scope-snappy is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
 * details.
 *
 * You should have received a copy of the GNU General Public License along with
 * unity-scope-snappy. If not, see <http://www.gnu.org/licenses/>.
 */

package daemon

import (
	"fmt"
	"github.com/godbus/dbus"
	"time"
)

// SnapdPackageManagerInterface implements a DBus interface for managing
// packages in snapd.
type SnapdPackageManagerInterface struct {
	dbusConnection DbusWrapper
	operationId    uint64

	pollPeriod time.Duration

	baseObjectPath dbus.ObjectPath

	progressSignalName string
	finishedSignalName string
	errorSignalName    string
}

// SnapdPackageManagerInterface creates a new SnapdPackageManagerInterface.
//
// Parameters:
// dbusConnection: Connection to the dbus bus.
// interfaceName: DBus interface name to implement.
// baseObjectPath: Base object path to use for signals.
//
// Returns:
// - New SnapdPackageManagerInterface
// - Error (nil if none)
func NewSnapdPackageManagerInterface(dbusConnection DbusWrapper,
	interfaceName string,
	baseObjectPath dbus.ObjectPath) (*SnapdPackageManagerInterface, error) {
	manager := &SnapdPackageManagerInterface{dbusConnection: dbusConnection}

	if !baseObjectPath.IsValid() {
		return nil, fmt.Errorf(`Invalid base object path: "%s"`, baseObjectPath)
	}

	manager.pollPeriod = time.Second

	manager.baseObjectPath = baseObjectPath

	manager.progressSignalName = interfaceName + ".progress"
	manager.finishedSignalName = interfaceName + ".finished"
	manager.errorSignalName = interfaceName + ".error"

	return manager, nil
}

// Install requests that snapd begin installation of a specific package, and
// then begins a polling job to provide progress feedback via the dbus
// connection.
//
// Parameters:
// packageId: ID of the package to be installed by snapd.
//
// Returns:
// - Object path over which the progress feedback will be provided.
// - DBus error (nil if none)
func (manager *SnapdPackageManagerInterface) Install(packageId string) (dbus.ObjectPath, *dbus.Error) {
	return "", dbus.NewError("org.freedesktop.DBus.Error.Failed",
		[]interface{}{"Not yet implemented"})
}

// Uninstall requests that Snapd begin uninstallation of a specific package, and
// then begins a polling job to provide progress feedback via the dbus
// connection.
//
// Parameters:
// packageId: ID of the package to be uninstalled by snapd.
//
// Returns:
// - Object path over which the progress feedback will be provided.
// - DBus error (nil if none)
func (manager *SnapdPackageManagerInterface) Uninstall(packageId string) (dbus.ObjectPath, *dbus.Error) {
	return "", dbus.NewError("org.freedesktop.DBus.Error.Failed",
		[]interface{}{"Not yet implemented"})
}
