package scope

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/webdm"
)

// FakePackageManager is a fake implementation of the PackageManager interface,
// for use within tests.
type FakePackageManager struct {
	getInstalledPackagesCalled bool
	getStorePackagesCalled     bool
	queryCalled                bool
	installCalled              bool
	uninstallCalled            bool

	failToGetInstalledPackages bool
	failToGetStorePackages     bool
	failToQuery                bool
	failToInstall              bool
	failToUninstall            bool
}

func (packageManager *FakePackageManager) GetInstalledPackages() ([]webdm.Package, error) {
	packageManager.getInstalledPackagesCalled = true

	if packageManager.failToGetInstalledPackages {
		return nil, fmt.Errorf("Failed to get installed packages (at user request)")
	}

	packages := make([]webdm.Package, 1)
	packages[0] = webdm.Package{Id: "package1", Status: webdm.StatusInstalled}

	return packages, nil
}

func (packageManager *FakePackageManager) GetStorePackages() ([]webdm.Package, error) {
	packageManager.getStorePackagesCalled = true

	if packageManager.failToGetStorePackages {
		return nil, fmt.Errorf("Failed to get store packages (at user request)")
	}

	packages := make([]webdm.Package, 1)
	packages[0] = webdm.Package{Id: "package1", Status: webdm.StatusNotInstalled}

	return packages, nil
}

func (packageManager *FakePackageManager) Query(packageId string) (*webdm.Package, error) {
	packageManager.queryCalled = true

	if packageManager.failToQuery {
		return nil, fmt.Errorf("Failed to query (at user request)")
	}

	return &webdm.Package{Id: packageId, Status: webdm.StatusNotInstalled}, nil
}

func (packageManager *FakePackageManager) Install(packageId string) error {
	packageManager.installCalled = true

	if packageManager.failToInstall {
		return fmt.Errorf("Failed to install (at user request)")
	}

	return nil
}

func (packageManager *FakePackageManager) Uninstall(packageId string) error {
	packageManager.uninstallCalled = true

	if packageManager.failToUninstall {
		return fmt.Errorf("Failed to uninstall (at user request)")
	}

	return nil
}
