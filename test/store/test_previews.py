#!/usr/bin/env python3

import os
import sys
import unittest
import fixtures

from scope_harness.testing import ScopeHarnessTestCase
from scope_harness import (ScopeHarness,
                           Parameters,
                           DepartmentMatcher,
                           ChildDepartmentMatcher,
                           PreviewView,
                           PreviewMatcher,
                           PreviewColumnMatcher,
                           PreviewWidgetMatcher)

# Local imports
from test.fakes import FakeWebdmServer
from test.test_fixtures import ServerFixture
import test.store.package_management_tasks as tasks

class TestPreviews(ScopeHarnessTestCase, fixtures.TestWithFixtures):
	"""Test the different package previews available in the store scope.
	"""

	def setUp(self):
		"""Setup the test fixtures and run the store scope using the harness.
		"""

		server = ServerFixture(FakeWebdmServer)
		self.useFixture(server)

		os.environ["WEBDM_URL"] = server.url

		self.harness = ScopeHarness.new_from_scope_list(Parameters([
			"{}/../../store/snappy-store.ini".format(os.path.dirname(os.path.realpath(__file__)))
		]))
		self.view = self.harness.results_view
		self.view.active_scope = "snappy-store"
		self.view.search_query = ""

	def testStorePreviewLayout(self):
		"""Test the preview of a package that is not installed.
		"""

		self.view.browse_department("")

		preview = self.view.category("store_packages").result("snappy:package2.canonical").tap()
		self.assertIsInstance(preview, PreviewView)

		tasks.verifyStorePreview(self, preview, "package2", "Canonical",
		                         "http://icon2", "description2", "0.2",
		                         "124 kB")

	def testInstalledPreviewLayout(self):
		"""Test the preview of a package that is installed.
		"""

		self.view.browse_department("")

		preview = self.view.category("store_packages").result("snappy:package1.canonical").tap()
		self.assertIsInstance(preview, PreviewView)

		tasks.verifyInstalledPreview(self, preview, "package1", "Canonical",
		                             "http://icon1", "description1", "0.1",
		                             "124 kB")
