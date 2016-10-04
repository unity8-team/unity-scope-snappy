/* Copyright (C) 2015 Canonical Ltd.
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

package fakes

import (
	"fmt"
	"github.com/godbus/dbus"
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
