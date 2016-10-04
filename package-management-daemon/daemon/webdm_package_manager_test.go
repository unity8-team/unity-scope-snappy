package daemon

import (
	"github.com/godbus/dbus"
	"testing"
	"time"
)

// Test typical NewWebdmPackageManagerInterface usage.
func TestNewWebdmPackageManagerInterface(t *testing.T) {
	manager, err := NewWebdmPackageManagerInterface(new(FakeDbusServer), "foo", "/foo", "")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	if manager.packageManager == nil {
		t.Error("Package manager was unexpectedly nil")
	}

	if manager.progressSignalName != "foo.progress" {
		t.Errorf(`Progress signal name was "%s", expected "foo.progress"`, manager.progressSignalName)
	}

	if manager.finishedSignalName != "foo.finished" {
		t.Errorf(`Finished signal name was "%s", expected "foo.finished"`, manager.finishedSignalName)
	}

	if manager.errorSignalName != "foo.error" {
		t.Errorf(`Error signal name was "%s", expected "foo.error"`, manager.errorSignalName)
	}
}

// Test that NewWebdmPackageManagerInterface fails with an invalid base object
// path.
func TestNewWebdmPackageManagerInterface_invalidBaseObjectPath(t *testing.T) {
	_, err := NewWebdmPackageManagerInterface(new(FakeDbusServer), "foo", "invalid", "")
	if err == nil {
		t.Error("Expected an error due to an invalid base object path")
	}
}

// Test that NewWebdmPackageManagerInterface fails with an invalid URL.
func TestNewWebdmPackageManagerInterface_invalidUrl(t *testing.T) {
	_, err := NewWebdmPackageManagerInterface(new(FakeDbusServer), "foo", "/foo", ":")
	if err == nil {
		t.Error("Expected an error due to an invalid API URL")
	}
}

// Test typical Install usage.
func TestInstall(t *testing.T) {
	dbusServer := new(FakeDbusServer)
	dbusServer.InitializeSignals()

	manager, err := NewWebdmPackageManagerInterface(dbusServer, "foo", "/foo", "")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	// Make the manager poll faster so the tests are more timely
	manager.pollPeriod = time.Millisecond

	packageManager := new(FakePackageManager)

	manager.packageManager = packageManager

	// Begin installation of two packages
	replyFoo, dbusErr := manager.Install("foo")
	if dbusErr != nil {
		t.Errorf(`Unexpected error while installing "foo": %s`, dbusErr)
	}

	replyBar, dbusErr := manager.Install("bar")
	if dbusErr != nil {
		t.Errorf(`Unexpected error while installing "bar": %s`, dbusErr)
	}

	if !packageManager.installCalled {
		t.Error("Expected package manager's Install method to be called!")
	}

	currentFooProgress := float32(0)
	currentBarProgress := float32(0)
	for signal := range dbusServer.signals {
		switch signal.Path {
		case replyFoo:
			if verifyFeedbackSignal(t, manager, signal, &currentFooProgress) {
				return
			}
		case replyBar:
			if verifyFeedbackSignal(t, manager, signal, &currentBarProgress) {
				return
			}
		default:
			t.Fatalf(`Signal path was "%s", expected either "%s" or "%s"`, signal.Path, replyFoo, replyBar)
		}
	}
}

// Test that failure during Install results in an error
func TestInstall_failure(t *testing.T) {
	manager, err := NewWebdmPackageManagerInterface(new(FakeDbusServer), "foo", "/foo", "")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	manager.packageManager = &FakePackageManager{failInstall: true}

	_, dbusErr := manager.Install("foo")
	if dbusErr == nil {
		t.Error("Expected error due to installation failure")
	}
}

// Data for TestInstall_inProgressFailure
var inProgressInstallFailureTests = []*FakePackageManager{
	&FakePackageManager{
		failInProgressInstall: true,
		failWithMessage:       true,
	},
	&FakePackageManager{
		failInProgressInstall: true,
		failWithMessage:       false,
	},
}

// Test that failure during in-progress Install results in an error, even
// if WebDM doesn't given a reason for the error.
func TestInstall_inProgressFailure(t *testing.T) {
	for i, packageManager := range inProgressInstallFailureTests {
		dbusServer := new(FakeDbusServer)
		dbusServer.InitializeSignals()

		manager, err := NewWebdmPackageManagerInterface(dbusServer, "foo", "/foo", "")
		if err != nil {
			t.Errorf("Test case %d: Unexpected error while creating new manager: %s", i, err)
			continue
		}

		manager.packageManager = packageManager

		reply, dbusErr := manager.Install("foo")
		if dbusErr != nil {
			t.Errorf(`Test case %d: Unexpected error while installing "foo": %s`, i, dbusErr)
		}

		signal := <-dbusServer.signals

		if signal.Path != reply {
			t.Fatalf(`Test case %d: Signal path was "%s", expected "%s"`, i, signal.Path, reply)
		}

		if signal.Name != manager.errorSignalName {
			t.Fatalf(`Test case %d: Signal name was "%s", expected "%s"`, i, signal.Name, manager.errorSignalName)
		}
	}
}

// Test that a Query error during installation results in an error being
// emitted.
func TestInstall_queryFailure(t *testing.T) {
	dbusServer := new(FakeDbusServer)
	dbusServer.InitializeSignals()

	manager, err := NewWebdmPackageManagerInterface(dbusServer, "foo", "/foo", "")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	manager.packageManager = &FakePackageManager{failQuery: true}

	reply, dbusErr := manager.Install("foo")
	if dbusErr != nil {
		t.Errorf(`Unexpected error while installing "foo": %s`, dbusErr)
	}

	signal := <-dbusServer.signals

	if signal.Path != reply {
		t.Fatalf(`Signal path was "%s", expected "%s"`, signal.Path, reply)
	}

	if signal.Name != manager.errorSignalName {
		t.Fatalf(`Signal name was "%s", expected "%s"`, signal.Name, manager.errorSignalName)
	}
}

// Test typical Uninstall usage.
func TestUninstall(t *testing.T) {
	dbusServer := new(FakeDbusServer)
	dbusServer.InitializeSignals()

	manager, err := NewWebdmPackageManagerInterface(dbusServer, "foo", "/foo", "")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	// Make the manager poll faster so the tests are more timely
	manager.pollPeriod = time.Millisecond

	packageManager := new(FakePackageManager)

	manager.packageManager = packageManager

	// Begin uninstallation of two packages
	replyFoo, dbusErr := manager.Uninstall("foo")
	if dbusErr != nil {
		t.Errorf(`Unexpected error while uninstalling "foo": %s`, dbusErr)
	}

	replyBar, dbusErr := manager.Uninstall("bar")
	if dbusErr != nil {
		t.Errorf(`Unexpected error while uninstalling "bar": %s`, dbusErr)
	}

	if !packageManager.uninstallCalled {
		t.Error("Expected package manager's Uninstall method to be called!")
	}

	currentFooProgress := float32(0)
	currentBarProgress := float32(0)
	for signal := range dbusServer.signals {
		switch signal.Path {
		case replyFoo:
			if verifyFeedbackSignal(t, manager, signal, &currentFooProgress) {
				return
			}
		case replyBar:
			if verifyFeedbackSignal(t, manager, signal, &currentBarProgress) {
				return
			}
		default:
			t.Fatalf(`Signal path was "%s", expected either "%s" or "%s"`, signal.Path, replyFoo, replyBar)
		}
	}
}

// Test that failure during Uninstall results in an error
func TestUninstall_failure(t *testing.T) {
	manager, err := NewWebdmPackageManagerInterface(new(FakeDbusServer), "foo", "/foo", "")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	manager.packageManager = &FakePackageManager{failUninstall: true}

	_, dbusErr := manager.Uninstall("foo")
	if dbusErr == nil {
		t.Error("Expected error due to uninstallation failure")
	}
}

// Data for TestUninstall_inProgressFailure
var inProgressUninstallFailureTests = []*FakePackageManager{
	&FakePackageManager{
		failInProgressUninstall: true,
		failWithMessage:         true,
	},
	&FakePackageManager{
		failInProgressUninstall: true,
		failWithMessage:         false,
	},
}

// Test that failure during in-progress Uninstall results in an error, even
// if WebDM doesn't given a reason for the error.
func TestUninstall_inProgressFailure(t *testing.T) {
	for i, packageManager := range inProgressUninstallFailureTests {
		dbusServer := new(FakeDbusServer)
		dbusServer.InitializeSignals()

		manager, err := NewWebdmPackageManagerInterface(dbusServer, "foo", "/foo", "")
		if err != nil {
			t.Errorf("Test case %d: Unexpected error while creating new manager: %s", i, err)
			continue
		}

		manager.packageManager = packageManager

		reply, dbusErr := manager.Uninstall("foo")
		if dbusErr != nil {
			t.Errorf(`Test case %d: Unexpected error while uninstalling "foo": %s`, i, dbusErr)
		}

		signal := <-dbusServer.signals

		if signal.Path != reply {
			t.Fatalf(`Test case %d: Signal path was "%s", expected "%s"`, i, signal.Path, reply)
		}

		if signal.Name != manager.errorSignalName {
			t.Fatalf(`Test case %d: Signal name was "%s", expected "%s"`, i, signal.Name, manager.errorSignalName)
		}
	}
}

// Test that a Query error during uninstallation results in an error being
// emitted.
func TestUninstall_queryFailure(t *testing.T) {
	dbusServer := new(FakeDbusServer)
	dbusServer.InitializeSignals()

	manager, err := NewWebdmPackageManagerInterface(dbusServer, "foo", "/foo", "")
	if err != nil {
		t.Fatalf("Unexpected error while creating new manager: %s", err)
	}

	manager.packageManager = &FakePackageManager{failQuery: true}

	reply, dbusErr := manager.Uninstall("foo")
	if dbusErr != nil {
		t.Errorf(`Unexpected error while uninstalling "foo": %s`, dbusErr)
	}

	signal := <-dbusServer.signals

	if signal.Path != reply {
		t.Fatalf(`Signal path was "%s", expected "%s"`, signal.Path, reply)
	}

	if signal.Name != manager.errorSignalName {
		t.Fatalf(`Signal name was "%s", expected "%s"`, signal.Name, manager.errorSignalName)
	}
}

func verifyFeedbackSignal(t *testing.T, manager *WebdmPackageManagerInterface, signal *dbus.Signal, currentProgress *float32) bool {
	switch signal.Name {
	case manager.progressSignalName:
		if len(signal.Body) != 2 {
			t.Fatalf("Got %d values, expected 2", len(signal.Body))
		}

		progress := signal.Body[0].(uint64)
		total := signal.Body[1].(uint64)

		if progress > total {
			t.Fatal("Progress is unexpectedly over the total. `finished` should have been emitted by now")
		}

		previousProgress := *currentProgress
		*currentProgress = float32(progress) / float32(total)

		if *currentProgress < previousProgress {
			t.Fatal("Installation isn't progressing as expected")
		}
	case manager.finishedSignalName:
		return true
	default:
		t.Fatalf("Unexpected signal name: %s", signal.Name)
	}

	return false
}
