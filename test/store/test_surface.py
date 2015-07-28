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
                           CategoryMatcher,
                           CategoryMatcherMode,
                           CategoryListMatcher,
                           CategoryListMatcherMode,
                           ResultMatcher)

# Local imports
from test.fakes import FakeWebdmServer
from test.test_fixtures import ServerFixture

class TestInitialSurface(ScopeHarnessTestCase, fixtures.TestWithFixtures):
	"""Test the initial surfacing of the store scope (no search query).
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

	def testScopeProperties(self):
		"""Test the scope properties from the .ini file"""

		self.assertEqual(self.view.scope_id, "snappy-store")
		self.assertEqual(self.view.display_name, "Snappy Scope")
		self.assertEqual(self.view.description, "A scope for the snappy store")

	def testStorePackages(self):
		"""Test that the store packages category contains the correct packages.
		"""

		self.assertMatchResult(CategoryListMatcher()
			.has_exactly(1)
			.mode(CategoryListMatcherMode.BY_ID)
			.category(CategoryMatcher("store_packages")
				.mode(CategoryMatcherMode.ALL)
				.result(ResultMatcher("snappy:package1.canonical")
					.title("package1")
					.art("http://icon1")
					.properties({"subtitle": "Canonical"})
					.properties({"id": "package1.canonical"}))
				.result(ResultMatcher("snappy:package2.canonical")
					.title("package2")
					.art("http://icon2")
					.properties({"subtitle": "Canonical"})
					.properties({"id": "package2.canonical"})))
		.match(self.view.categories))

	def testInstalledPackages(self):
		"""Test that the installed packages category contains only the
		   installed package.
		"""

		self.view.browse_department("installed")

		self.assertMatchResult(CategoryListMatcher()
			.has_exactly(1)
			.mode(CategoryListMatcherMode.BY_ID)
			.category(CategoryMatcher("installed_packages")
				.mode(CategoryMatcherMode.ALL)
				.result(ResultMatcher("snappy:package1.canonical")
					.title("package1")
					.art("http://icon1")
					.properties({"subtitle": "Canonical"})
					.properties({"id": "package1.canonical"})))
		.match(self.view.categories))
