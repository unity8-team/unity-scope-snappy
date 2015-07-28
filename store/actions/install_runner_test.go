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

package actions

import (
	"launchpad.net/unity-scope-snappy/internal/github.com/godbus/dbus"
	"launchpad.net/unity-scope-snappy/internal/launchpad.net/go-unityscopes/v2"
	"launchpad.net/unity-scope-snappy/store/operation"
	"launchpad.net/unity-scope-snappy/store/packages/fakes"
	"testing"
)

// Test typical Run usage.
func TestInstallRunner_run(t *testing.T) {
	actionRunner, _ := NewInstallRunner()

	packageManager := new(fakes.FakeDbusManager)

	response, err := actionRunner.Run(packageManager, "foo")
	if err != nil {
		// Exit here so we don't dereference nil
		t.Fatalf("Unexpected error when attempting to run: %s", err)
	}

	if !packageManager.InstallCalled {
		t.Error("Expected package manager Install() function to be called")
	}

	if response.Status != scopes.ActivationShowPreview {
		t.Errorf(`Response status was "%d", expected "%d"`, response.Status, scopes.ActivationShowPreview)
	}

	// Verify operation metadata
	metadata, ok := response.ScopeData.(operation.Metadata)
	if !ok {
		// Exit here so we don't dereference nil
		t.Fatalf("Expected response ScopeData to include operation metadata")
	}

	if !metadata.InstallRequested {
		t.Errorf("Expected metadata to indicate that an installation was requested")
	}

	if metadata.ObjectPath != dbus.ObjectPath("/foo/1") {
		t.Errorf(`Metadata object path was "%s", expected "/foo/1"`, metadata.ObjectPath)
	}
}

// Test that a failure to install results in an error
func TestInstallRunner_run_installationFailure(t *testing.T) {
	actionRunner, _ := NewInstallRunner()

	packageManager := &fakes.FakeDbusManager{FailInstall: true}

	response, err := actionRunner.Run(packageManager, "foo")
	if err == nil {
		t.Error("Expected an error due to failure to install")
	}
	if response != nil {
		t.Error("Expected response to be nil")
	}
}
