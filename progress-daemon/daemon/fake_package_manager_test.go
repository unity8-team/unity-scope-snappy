package daemon

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/webdm"
	"testing"
)

const (
	progressStep = 50
)

// FakePackageManager is a fake implementation of the PackageManager interface,
// for use within tests.
type FakePackageManager struct {
	queryCalled     bool
	installCalled   bool
	uninstallCalled bool

	failQuery     bool
	failInstall   bool
	failUninstall bool

	failInProgressInstall   bool
	failInProgressUninstall bool

	failWithMessage bool

	// Key: Package ID
	// Value: Install progress (0-100)
	installingPackages map[string]float64

	// Key: Package ID
	// Value: Install progress (0-100)
	uninstallingPackages map[string]float64
}

func (packageManager *FakePackageManager) Query(packageId string) (*webdm.Package, error) {
	packageManager.queryCalled = true

	if packageManager.failQuery {
		return nil, fmt.Errorf("Failed at user request")
	}

	snap := &webdm.Package{Id: packageId, Status: webdm.StatusNotInstalled}

	if packageManager.installingPackages != nil {
		progress, ok := packageManager.installingPackages[packageId]
		if ok {
			progress = continueOperation(progress, snap,
				webdm.StatusInstalling, webdm.StatusInstalled,
				webdm.StatusNotInstalled, packageManager.failInProgressInstall,
				packageManager.failWithMessage)
			packageManager.installingPackages[packageId] = progress
		}
	}

	if packageManager.uninstallingPackages != nil {
		progress, ok := packageManager.uninstallingPackages[packageId]
		if ok {
			progress = continueOperation(progress, snap,
				webdm.StatusUninstalling, webdm.StatusNotInstalled,
				webdm.StatusInstalled, packageManager.failInProgressUninstall,
				packageManager.failWithMessage)
			packageManager.uninstallingPackages[packageId] = progress
		}
	}

	return snap, nil
}

func (packageManager *FakePackageManager) Install(packageId string) error {
	packageManager.installCalled = true

	if packageManager.failInstall {
		return fmt.Errorf("Failed at user request")
	}

	if packageManager.installingPackages == nil {
		packageManager.installingPackages = make(map[string]float64)
	}

	// Set install progress to 0%
	packageManager.installingPackages[packageId] = 0

	return nil
}

func (packageManager *FakePackageManager) Uninstall(packageId string) error {
	packageManager.uninstallCalled = true

	if packageManager.failUninstall {
		return fmt.Errorf("Failed at user request")
	}

	if packageManager.uninstallingPackages == nil {
		packageManager.uninstallingPackages = make(map[string]float64)
	}

	// Set uninstall progress to 0%
	packageManager.uninstallingPackages[packageId] = 0

	return nil
}

func continueOperation(progress float64, snap *webdm.Package,
	inProgressStatus webdm.Status, finishedStatus webdm.Status,
	errorStatus webdm.Status, fail bool, failWithMessage bool) float64 {
	if fail {
		snap.Status = errorStatus

		if failWithMessage {
			snap.Message = "Failed at user request"
		}

		return 0.0
	}

	if progress < 100 {
		// Operation isn't "done" yet. Keep going.
		snap.Status = inProgressStatus
		progress += progressStep
		snap.Progress = progress
	} else {
		snap.Status = finishedStatus
	}

	return progress
}

// Test that an Install followed by a Query shows install progress as expected
func TestFakePackageManager_installProgress(t *testing.T) {
	packageManager := new(FakePackageManager)

	err := packageManager.Install("foo")
	if err != nil {
		t.Errorf("Unexpected error when installing: %s", err)
	}

	for i := 1; i <= 100/progressStep; i++ {
		snap, err := packageManager.Query("foo")
		if err != nil {
			t.Errorf("Unexpected error when querying: %s", err)
		}

		if snap.Status != webdm.StatusInstalling {
			t.Errorf("Status was %d, expected %d", snap.Status, webdm.StatusInstalling)
		}

		expected := float64(progressStep * i)

		if snap.Progress != expected {
			t.Errorf("Progress was %f, expected %f", snap.Progress, expected)
		}
	}
}

// Test that an Uninstall followed by a Query shows uninstall progress as
// expected
func TestFakePackageManager_uninstallProgress(t *testing.T) {
	packageManager := new(FakePackageManager)

	err := packageManager.Uninstall("foo")
	if err != nil {
		t.Errorf("Unexpected error when uninstalling: %s", err)
	}

	for i := 1; i <= 100/progressStep; i++ {
		snap, err := packageManager.Query("foo")
		if err != nil {
			t.Errorf("Unexpected error when querying: %s", err)
		}

		if snap.Status != webdm.StatusUninstalling {
			t.Errorf("Status was %d, expected %d", snap.Status, webdm.StatusUninstalling)
		}

		expected := float64(progressStep * i)

		if snap.Progress != expected {
			t.Errorf("Progress was %f, expected %f", snap.Progress, expected)
		}
	}
}
