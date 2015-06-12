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
	installingPackages map[string]int

	// Key: Package ID
	// Value: Install progress (0-100)
	uninstallingPackages map[string]int
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
			if packageManager.failInProgressInstall {
				snap.Status = webdm.StatusNotInstalled

				if packageManager.failWithMessage {
					snap.Message = "Failed at user request"
				}
			} else {
				if progress < 100 {
					// Package isn't installed yet. Keep "installing" it.
					snap.Status = webdm.StatusInstalling
					progress += progressStep
					snap.Progress = float64(progress)
					packageManager.installingPackages[packageId] = progress
				} else {
					snap.Status = webdm.StatusInstalled
				}
			}
		}
	}

	if packageManager.uninstallingPackages != nil {
		progress, ok := packageManager.uninstallingPackages[packageId]
		if ok {
			if packageManager.failInProgressUninstall {
				snap.Status = webdm.StatusInstalled

				if packageManager.failWithMessage {
					snap.Message = "Failed at user request"
				}
			} else {
				if progress < 100 {
					// Package isn't installed yet. Keep "installing" it.
					snap.Status = webdm.StatusUninstalling
					progress += 50
					snap.Progress = float64(progress)
					packageManager.uninstallingPackages[packageId] = progress
				} else {
					snap.Status = webdm.StatusNotInstalled
				}
			}
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
		packageManager.installingPackages = make(map[string]int)
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
		packageManager.uninstallingPackages = make(map[string]int)
	}

	// Set uninstall progress to 0%
	packageManager.uninstallingPackages[packageId] = 0

	return nil
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

		expected := progressStep * i
		actual := int(snap.Progress)

		if actual != expected {
			t.Errorf("Progress was %d, expected %d", actual, expected)
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

		expected := progressStep * i
		actual := int(snap.Progress)

		if actual != expected {
			t.Errorf("Progress was %d, expected %d", actual, expected)
		}
	}
}
