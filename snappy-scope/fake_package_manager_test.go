package main

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/webdm"
)

// FakePackageManager is a fake implementation of the PackageManager interface,
// for use within tests.
type FakePackageManager struct {
	installCalled   bool
	uninstallCalled bool
}

func (packageManager FakePackageManager) GetInstalledPackages() ([]webdm.Package, error) {
	return nil, fmt.Errorf("Not implemented...")
}

func (packageManager FakePackageManager) GetStorePackages() ([]webdm.Package, error) {
	return nil, fmt.Errorf("Not implemented...")
}

func (packageManager FakePackageManager) Query(packageId string) (*webdm.Package, error) {
	return nil, fmt.Errorf("Not implemented...")
}

func (packageManager *FakePackageManager) Install(packageId string) error {
	packageManager.installCalled = true
	return nil
}

func (packageManager *FakePackageManager) Uninstall(packageId string) error {
	packageManager.uninstallCalled = true
	return nil
}
