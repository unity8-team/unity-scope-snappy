#!/usr/bin/env python3

# Copyright (C) 2015 Canonical Ltd.
#
# This file is part of unity-scope-snappy.
#
# unity-scope-snappy is free software: you can redistribute it and/or modify it
# under the terms of the GNU General Public License as published by the Free
# Software Foundation, either version 3 of the License, or (at your option) any
# later version.
#
# unity-scope-snappy is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
# FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more
# details.
#
# You should have received a copy of the GNU General Public License along with
# unity-scope-snappy. If not, see <http://www.gnu.org/licenses/>.

import os
import sys
import subprocess
import unittest
import fixtures

from scope_harness.testing import ScopeHarnessTestCase
from scope_harness import (ScopeHarness,
                           Parameters,
                           PreviewView,
                           PreviewColumnMatcher,
                           PreviewMatcher,
                           PreviewWidgetMatcher)

# Local imports
from test.fakes import FakeWebdmServer
from test.test_fixtures import ServerFixture
import test.store.package_management_tasks as tasks

THIS_FILE_PATH = os.path.dirname(os.path.realpath(__file__))

class TestPackageManagement(ScopeHarnessTestCase, fixtures.TestWithFixtures):
	"""Test package installation/uninstallation via the scope.
	"""

	def setUp(self):
		"""Setup the test fixtures, run the store scope, and fire up the daemon.
		"""

		server = ServerFixture(FakeWebdmServer)
		self.useFixture(server)

		# Run the package management daemon to communicate with our fake server
		self.daemon = subprocess.Popen(["package-management-daemon",
		                                "-webdm={}".format(server.url)])

		os.environ["WEBDM_URL"] = server.url

		self.harness = ScopeHarness.new_from_scope_list(Parameters([
			"{}/../../store/snappy-store.ini".format(THIS_FILE_PATH)
		]))
		self.view = self.harness.results_view
		self.view.active_scope = "snappy-store"
		self.view.search_query = ""

	def tearDown(self):
		"""Terminate the daemon and make sure it exits cleanly."""

		self.daemon.terminate()

		# Make sure the daemon exits cleanly
		self.assertEqual(self.daemon.wait(), 0)

	def testPackageInstallation(self):
		"""Test installing a package via the scope."""

		self.installPackage("store_packages", "package2.canonical", "package2",
		                    "Canonical", "http://icon2", "description2", "0.2",
		                    "124 kB")

	def testPackageUninstallation(self):
		"""Test uninstalling a package via the scope."""

		self.uninstallPackage("store_packages", "package1.canonical",
		                      "package1", "Canonical", "http://icon1",
		                      "description1", "0.1", "124 kB")

	def testPackageInstallUninstallReinstall(self):
		"""Test handling of package install/uninstall/reinstall flow."""

		# Install the package
		self.installPackage("store_packages", "package2.canonical", "package2",
		                    "Canonical", "http://icon2", "description2", "0.2",
		                    "124 kB")

		# Now unstall the package
		self.uninstallPackage("store_packages", "package2.canonical",
		                      "package2", "Canonical", "http://icon2",
		                      "description2", "0.2", "124 kB", 2)

		# Finally, re-install the package
		self.installPackage("store_packages", "package2.canonical", "package2",
		                    "Canonical", "http://icon2", "description2", "0.2",
		                    "124 kB", 3)

	def installPackage(self, categoryId, packageId, title, subtitle, mascot,
	                   description, version, size, operationNumber=1):
		"""Install and verify package installation via scope."""

		preview = tasks.installPackage(self, categoryId, packageId)
		self.assertIsInstance(preview, PreviewView)

		tasks.verifyInstallingPreview(self, preview, title, subtitle, mascot,
		                              description, version, size,
		                              operationNumber)

		# Force the installation to succeed
		widgetFound = False
		for columnWidgets in preview.widgets:
			for widget in columnWidgets:
				if widget.type == "progress":
					preview = widget.trigger("finished", "")
					widgetFound = True

		self.assertTrue(widgetFound,
		                "Expected progress widget to be in preview.")

		tasks.verifyInstalledPreview(self, preview, title, subtitle, mascot,
		                             description, version, size)

	def uninstallPackage(self, categoryId, packageId, title, subtitle, mascot,
	                     description, version, size, operationNumber=1):
		"""Uninstall and verify package uninstallation via scope."""

		preview = tasks.uninstallPackage(self, categoryId, packageId)
		self.assertIsInstance(preview, PreviewView)

		tasks.verifyConfirmUninstallPreview(self, preview, title)

		preview = tasks.confirmUninstallPackage(self, preview)
		self.assertIsInstance(preview, PreviewView)

		tasks.verifyUninstallingPreview(self, preview, title, subtitle, mascot,
		                                description, version, size,
		                                operationNumber)

		# Force the uninstallation to succeed
		widgetFound = False
		for columnWidgets in preview.widgets:
			for widget in columnWidgets:
				if widget.type == "progress":
					preview = widget.trigger("finished", "")
					widgetFound = True

		self.assertTrue(widgetFound,
		                "Expected progress widget to be in preview.")

		tasks.verifyStorePreview(self, preview, title, subtitle, mascot,
		                         description, version, size)
