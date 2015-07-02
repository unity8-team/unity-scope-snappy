package fakes

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/webdm"
)

// FakeManager is a fake implementation of the Manager interface, for use within
// tests.
type FakeManager struct {
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

func (manager *FakeManager) GetInstalledPackages() ([]webdm.Package, error) {
	manager.GetInstalledPackagesCalled = true

	if manager.FailGetInstalledPackages {
		return nil, fmt.Errorf("Failed to get installed packages (at user request)")
	}

	packages := make([]webdm.Package, 1)
	packages[0] = webdm.Package{Id: "package1", Status: webdm.StatusInstalled}

	return packages, nil
}

func (manager *FakeManager) GetStorePackages() ([]webdm.Package, error) {
	manager.GetStorePackagesCalled = true

	if manager.FailGetStorePackages {
		return nil, fmt.Errorf("Failed to get store packages (at user request)")
	}

	packages := make([]webdm.Package, 1)
	packages[0] = webdm.Package{Id: "package1", Status: webdm.StatusNotInstalled}

	return packages, nil
}

func (manager *FakeManager) Query(packageId string) (*webdm.Package, error) {
	manager.QueryCalled = true

	if manager.FailQuery {
		return nil, fmt.Errorf("Failed to query (at user request)")
	}

	return &webdm.Package{Id: packageId, Status: webdm.StatusNotInstalled}, nil
}

func (manager *FakeManager) Install(packageId string) error {
	manager.InstallCalled = true

	if manager.FailInstall {
		return fmt.Errorf("Failed to install (at user request)")
	}

	return nil
}

func (manager *FakeManager) Uninstall(packageId string) error {
	manager.UninstallCalled = true

	if manager.FailUninstall {
		return fmt.Errorf("Failed to uninstall (at user request)")
	}

	return nil
}
