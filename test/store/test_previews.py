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
