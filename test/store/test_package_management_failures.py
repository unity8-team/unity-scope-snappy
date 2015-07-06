#!/usr/bin/env python3

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

class TestPackageManagementFailures(ScopeHarnessTestCase, fixtures.TestWithFixtures):
	"""Test package installation/uninstallation failures in the scope.
	"""

	def setUp(self):
		"""Setup the test fixtures, run the store scope, and fire up the daemon.
		"""

		server = ServerFixture(FakeWebdmServer, True)
		self.useFixture(server)

		# Run the progress daemon to communicate with our fake server
		self.daemon = subprocess.Popen(["progress-daemon",
		                                "-webdm={}".format(server.url)])

		os.environ["WEBDM_URL"] = server.url

		self.harness = ScopeHarness.new_from_scope_list(Parameters([
			"{}/../../scope.ini".format(THIS_FILE_PATH)
		]))
		self.view = self.harness.results_view
		self.view.active_scope = "scope"
		self.view.search_query = ""

	def tearDown(self):
		"""Terminate the daemon and make sure it exits cleanly."""

		self.daemon.terminate()

		# Make sure the daemon exits cleanly
		self.assertEqual(self.daemon.wait(), 0)

	def testPackageInstallationFailure(self):
		"""Test failure handling from the progress widget while installing."""

		preview = tasks.installPackage(self, "store_packages", "package2.canonical")
		self.assertIsInstance(preview, PreviewView)

		# Force the installation to fail
		widgetFound = False
		for columnWidgets in preview.widgets:
			for widget in columnWidgets:
				if widget.type == "progress":
					preview = widget.trigger("failed", "")
					widgetFound = True

		self.assertTrue(widgetFound,
		                "Expected progress widget to be in preview.")

		# Verify that the package still shows the package as not installed.
		tasks.verifyStorePreview(self, preview, "package2", "Canonical",
		                        "http://icon2", "description2", "0.2", "124 kB")

	def testPackageUninstallationFailure(self):
		"""Test failure handling from the progress widget while uninstalling."""

		preview = tasks.uninstallPackage(self, "store_packages", "package1.canonical")
		self.assertIsInstance(preview, PreviewView)

		tasks.verifyConfirmUninstallPreview(self, preview, "package1")

		preview = tasks.confirmUninstallPackage(self, preview)
		self.assertIsInstance(preview, PreviewView)

		# Force the uninstallation to fail
		widgetFound = False
		for columnWidgets in preview.widgets:
			for widget in columnWidgets:
				if widget.type == "progress":
					preview = widget.trigger("failed", "")
					widgetFound = True

		self.assertTrue(widgetFound,
		                "Expected progress widget to be in preview.")

		# Verify that the preview still shows the package as installed.
		tasks.verifyInstalledPreview(self, preview, "package1", "Canonical",
		                             "http://icon1", "description1", "0.1",
		                             "124 kB")
