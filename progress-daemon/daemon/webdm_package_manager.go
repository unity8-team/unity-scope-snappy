package daemon

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
	"launchpad.net/unity-scope-snappy/webdm"
	"time"
)

const (
	webdmDefaultApiUrl = webdm.DefaultApiUrl
)

// WebdmPackageManagerInterface implements a DBus interface for managing
// packages in WebDM.
type WebdmPackageManagerInterface struct {
	dbusConnection DbusWrapper
	packageManager PackageManager
	operationId    uint64

	pollPeriod time.Duration

	baseObjectPath dbus.ObjectPath

	progressSignalName string
	finishedSignalName string
	errorSignalName    string
}

// NewWebdmPackageManagerInterface creates a new WebdmPackageManagerInterface.
//
// Parameters:
// dbusConnection: Connection to the dbus bus.
// interfaceName: DBus interface name to implement.
// baseObjectPath: Base object path to use for signals.
// apiUrl: WebDM API URL.
//
// Returns:
// - New WebdmPackageManagerInterface
// - Error (nil if none)
func NewWebdmPackageManagerInterface(dbusConnection DbusWrapper,
	interfaceName string, baseObjectPath dbus.ObjectPath,
	apiUrl string) (*WebdmPackageManagerInterface, error) {
	manager := &WebdmPackageManagerInterface{dbusConnection: dbusConnection}

	if apiUrl == "" {
		apiUrl = webdmDefaultApiUrl
	}

	var err error
	manager.packageManager, err = webdm.NewClient(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to create webcm client: %s", err)
	}

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

// Install requests that WebDM begin installation of a specific package, and
// then begins a polling job to provide progress feedback via the dbus
// connection.
//
// Parameters:
// packageId: ID of the package to be installed by WebDM.
//
// Returns:
// - Object path over which the progress feedback will be provided.
// - DBus error (nil if none)
func (manager *WebdmPackageManagerInterface) Install(packageId string) (dbus.ObjectPath, *dbus.Error) {
	err := manager.packageManager.Install(packageId)
	if err != nil {
		return "", dbus.NewError("org.freedesktop.DBus.Error.Failed",
			[]interface{}{fmt.Sprintf(`Unable to install package "%s": %s`, packageId, err)})
	}

	operationId := manager.newOperationId()

	go manager.reportProgress(operationId, packageId, webdm.StatusInstalling, webdm.StatusInstalled)

	return manager.operationObjectPath(operationId), nil
}

// Uninstall requests that WebDM begin uninstallation of a specific package, and
// then begins a polling job to provide progress feedback via the dbus
// connection.
//
// Parameters:
// packageId: ID of the package to be uninstalled by WebDM.
//
// Returns:
// - Object path over which the progress feedback will be provided.
// - DBus error (nil if none)
func (manager *WebdmPackageManagerInterface) Uninstall(packageId string) (dbus.ObjectPath, *dbus.Error) {
	err := manager.packageManager.Uninstall(packageId)
	if err != nil {
		return "", dbus.NewError("org.freedesktop.DBus.Error.Failed",
			[]interface{}{fmt.Sprintf(`Unable to uninstall package "%s": %s`, packageId, err)})
	}

	operationId := manager.newOperationId()

	go manager.reportProgress(operationId, packageId, webdm.StatusUninstalling, webdm.StatusNotInstalled)

	return manager.operationObjectPath(operationId), nil
}

// reportProgress is the "polling job" used by both Install and Uninstall. It
// simply queries WebDM for the given package once a second, and fires off the
// `progress`, `finished`, or `error` DBus signals as appropriate.
//
// Parameters:
// packageId: ID of the package to be monitored.
// progressStatus: The expected status of the package while it's undergoing the
//                 anticipated action.
// finishedStatus: The status of the package when its finished the anticipated
//                 action.
func (manager *WebdmPackageManagerInterface) reportProgress(operationId uint64, packageId string, progressStatus webdm.Status, finishedStatus webdm.Status) {
	for {
		snap, err := manager.packageManager.Query(packageId)
		if err != nil {
			manager.emitError(operationId, `Unable to query package "%s": %s`, packageId, err)
			return
		}

		switch snap.Status {
		// If the package status still shows in-progress, report out the current
		// progress.
		case progressStatus:
			// Round the progress to the nearest integer
			progress := uint64(snap.Progress + .5)

			manager.emitProgress(operationId, progress, 100)

		// If the package status is what was desired, we can stop polling and
		// exit the progress loop, reporting success.
		case finishedStatus:
			manager.emitFinished(operationId)
			return

		// If the package status is anything other than in-progress or desired,
		// an error occurred and we can exit the progress loop, reporting the
		// failure.
		default:
			if snap.Message == "" {
				snap.Message = "(no message given)"
			}

			manager.emitError(operationId, `Failed to install package "%s": %s`, packageId, snap.Message)
			return
		}

		time.Sleep(manager.pollPeriod)
	}
}

// emitProgress emits the `progres` DBus signal.
//
// Parameters:
// packageId: ID of the package whose progress is being emitted.
// received: Received count.
// total: Total count.
func (manager *WebdmPackageManagerInterface) emitProgress(operationId uint64, received uint64, total uint64) {
	manager.dbusConnection.Emit(manager.operationObjectPath(operationId),
		manager.progressSignalName, received, total)
}

// emitFinished emits the `finished` DBus signal.
//
// Parameters:
// packageId: ID of the package that just finished an operation.
func (manager *WebdmPackageManagerInterface) emitFinished(operationId uint64) {
	manager.dbusConnection.Emit(manager.operationObjectPath(operationId),
		manager.finishedSignalName, "")
}

// emitError emits the `error` DBus signal.
//
// Parameters:
// packageId: ID of the package that encountered an error.
// format: Format string of the error.
// a...: List of values for the placeholders in the `format` string.
func (manager *WebdmPackageManagerInterface) emitError(operationId uint64, format string, a ...interface{}) {
	manager.dbusConnection.Emit(manager.operationObjectPath(operationId),
		manager.errorSignalName, fmt.Sprintf(format, a...))
}

// newOperationId is used to generate a unique ID to be used within dbus object
// paths. Normally we'd just use WebDM's package IDs to create unique dbus
// object paths, but package IDs can include characters that would be invalid
// within an object path. So we use this instead.
//
// Returns:
// - Unique number.
func (manager *WebdmPackageManagerInterface) newOperationId() uint64 {
	manager.operationId++
	return manager.operationId
}

// operationObjectPath is used to generate an object path for a given operation.
//
// Parameters:
// operationId: ID of the operation to be represented by this object path.
//
// Returns:
// - New object path.
func (manager WebdmPackageManagerInterface) operationObjectPath(operationId uint64) dbus.ObjectPath {
	return dbus.ObjectPath(fmt.Sprintf("%s/%d", manager.baseObjectPath, operationId))
}
