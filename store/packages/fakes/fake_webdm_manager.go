package fakes

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/webdm"
)

// FakeWebdmManager is a fake implementation of the WebdmManager interface, for
// use within tests.
type FakeWebdmManager struct {
	GetInstalledPackagesCalled bool
	GetStorePackagesCalled     bool
	QueryCalled                bool
	InstallCalled              bool
	UninstallCalled            bool

	FailGetInstalledPackages bool
	FailGetStorePackages     bool
	FailQuery                bool
	FailInstall              bool
	FailUninstall            bool
}

func (manager *FakeWebdmManager) GetInstalledPackages() ([]webdm.Package, error) {
	manager.GetInstalledPackagesCalled = true

	if manager.FailGetInstalledPackages {
		return nil, fmt.Errorf("Failed to get installed packages (at user request)")
	}

	packages := make([]webdm.Package, 1)
	packages[0] = webdm.Package{Id: "package1", Status: webdm.StatusInstalled}

	return packages, nil
}

func (manager *FakeWebdmManager) GetStorePackages() ([]webdm.Package, error) {
	manager.GetStorePackagesCalled = true

	if manager.FailGetStorePackages {
		return nil, fmt.Errorf("Failed to get store packages (at user request)")
	}

	packages := make([]webdm.Package, 1)
	packages[0] = webdm.Package{Id: "package1", Status: webdm.StatusNotInstalled}

	return packages, nil
}

func (manager *FakeWebdmManager) Query(packageId string) (*webdm.Package, error) {
	manager.QueryCalled = true

	if manager.FailQuery {
		return nil, fmt.Errorf("Failed to query (at user request)")
	}

	return &webdm.Package{Id: packageId, Status: webdm.StatusNotInstalled}, nil
}

func (manager *FakeWebdmManager) Install(packageId string) error {
	manager.InstallCalled = true

	if manager.FailInstall {
		return fmt.Errorf("Failed to install (at user request)")
	}

	return nil
}

func (manager *FakeWebdmManager) Uninstall(packageId string) error {
	manager.UninstallCalled = true

	if manager.FailUninstall {
		return fmt.Errorf("Failed to uninstall (at user request)")
	}

	return nil
}
