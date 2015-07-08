package fakes

import (
	"fmt"
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
)

// FakeDbusManager is a fake implementation of the DbusManager interface, for
// use within tests.
type FakeDbusManager struct {
	ConnectCalled   bool
	InstallCalled   bool
	UninstallCalled bool

	FailConnect   bool
	FailInstall   bool
	FailUninstall bool
}

func (manager *FakeDbusManager) Connect() error {
	manager.ConnectCalled = true

	if manager.FailConnect {
		return fmt.Errorf("Failed at user request")
	}

	return nil
}

func (manager *FakeDbusManager) Install(packageId string) (dbus.ObjectPath, error) {
	manager.InstallCalled = true

	if manager.FailInstall {
		return "", fmt.Errorf("Failed at user request")
	}

	return "/foo/1", nil
}

func (manager *FakeDbusManager) Uninstall(packageId string) (dbus.ObjectPath, error) {
	manager.UninstallCalled = true

	if manager.FailUninstall {
		return "", fmt.Errorf("Failed at user request")
	}

	return "/foo/1", nil
}
