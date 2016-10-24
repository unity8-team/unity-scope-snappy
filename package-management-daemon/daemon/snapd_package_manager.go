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
	"github.com/snapcore/snapd/client"
	"time"
)

// SnapdPackageManagerInterface implements a DBus interface for managing
// packages in snapd.
type SnapdPackageManagerInterface struct {
	dbusConnection DbusWrapper
	operationId    uint64

	pollPeriod time.Duration

	baseObjectPath dbus.ObjectPath

	clientConfig client.Config
	client *client.Client

	processingSignalName string
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

	manager.client = client.New(&manager.clientConfig)

	manager.processingSignalName = interfaceName + ".processing"
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
	opts := &client.SnapOptions{}

	var err error
	var changeID string

	changeID, err = manager.client.Install(packageId, opts)
	if err != nil {
		return "", dbus.NewError("org.freedesktop.DBus.Error.Failed",
			[]interface{}{fmt.Sprintf("Error installing package '%s': %s",
				packageId, err)})
	}

	go manager.wait(changeID)
	return manager.getObjectPath(changeID), nil
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
	opts := &client.SnapOptions{}

	var err error
	var changeID string

	changeID, err = manager.client.Remove(packageId, opts)
	if err != nil {
		return "", dbus.NewError("org.freedesktop.DBus.Error.Failed",
			[]interface{}{fmt.Sprintf("Error installing package '%s': %s",
				packageId, err)})
	}

	go manager.wait(changeID)
	return manager.getObjectPath(changeID), nil
}


// operationObjectPath is used to generate an object path for a given operation.
//
// Parameters:
// operationId: ID of the operation to be represented by this object path.
//
// Returns:
// - New object path.
func (manager *SnapdPackageManagerInterface) getObjectPath(changeID string) dbus.ObjectPath {
	return dbus.ObjectPath(fmt.Sprintf("%s/%s",
		manager.baseObjectPath, changeID))
}

// emitProgress emits the `progres` DBus signal.
//
// Parameters:
// packageId: ID of the package whose progress is being emitted.
// received: Received count.
// total: Total count.
func (manager *SnapdPackageManagerInterface) emitProgress(changeID string, received uint64, total uint64) {
	manager.dbusConnection.Emit(manager.getObjectPath(changeID),
		manager.progressSignalName, received, total)
}

// emitProcessing emits the `processing` DBus signal.
//
// Parameters:
// packageId: ID of the package that is processing.
func (manager *SnapdPackageManagerInterface) emitProcessing(changeID string) {
	manager.dbusConnection.Emit(manager.getObjectPath(changeID),
		manager.processingSignalName, "")
}

// emitFinished emits the `finished` DBus signal.
//
// Parameters:
// packageId: ID of the package that just finished an operation.
func (manager *SnapdPackageManagerInterface) emitFinished(changeID string) {
	manager.dbusConnection.Emit(manager.getObjectPath(changeID),
		manager.finishedSignalName, "")

	// Refresh the apps scope and the snappy store scope
	manager.dbusConnection.Emit("/com/canonical/unity/scopes",
		"com.canonical.unity.scopes.InvalidateResults",
		"clickscope")
	manager.dbusConnection.Emit("/com/canonical/unity/scopes",
		"com.canonical.unity.scopes.InvalidateResults",
		"snappy-store")
}

// emitError emits the `error` DBus signal.
//
// Parameters:
// packageId: ID of the package that encountered an error.
// format: Format string of the error.
// a...: List of values for the placeholders in the `format` string.
func (manager *SnapdPackageManagerInterface) emitError(changeID string, format string, a ...interface{}) {
	manager.dbusConnection.Emit(manager.getObjectPath(changeID),
		manager.errorSignalName, fmt.Sprintf(format, a...))
}

func (manager *SnapdPackageManagerInterface) wait(changeID string) {
	tMax := time.Time{}

	var lastID string
	for {
		chg, err := manager.client.Change(changeID)
		if err != nil {
			// an error here means the server most likely went away
			// XXX: it actually can be a bunch of other things; fix client to expose it better
			now := time.Now()
			if tMax.IsZero() {
				tMax = now.Add(manager.pollPeriod * 5)
			}
			if now.After(tMax) {
				manager.emitError(changeID, "Error talking to snapd: %s", err)
				return
			}
			manager.emitProcessing(changeID)
			time.Sleep(manager.pollPeriod)
			continue
		}
		if !tMax.IsZero() {
			tMax = time.Time{}
		}

		for _, t := range chg.Tasks {
			switch {
			case t.Status != "Doing":
				continue
			case t.Progress.Total == 1:
				manager.emitProcessing(changeID)
			case t.ID == lastID:
				manager.emitProgress(changeID, uint64(t.Progress.Done), uint64(t.Progress.Total))
			default:
				lastID = t.ID
			}
			break
		}

		if chg.Ready {
			if chg.Status == "Done" {
				manager.emitFinished(changeID)
			} else if chg.Err != "" {
				manager.emitError(changeID, chg.Err)
			}

			return
		}

		// note this very purposely is not a ticker; we want
		// to sleep 100ms between calls, not call once every
		// 100ms.
		time.Sleep(manager.pollPeriod)
	}
}
