#!/usr/bin/env python3

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
			"{}/../../scope.ini".format(os.path.dirname(os.path.realpath(__file__)))
		]))
		self.view = self.harness.results_view
		self.view.active_scope = "scope"
		self.view.search_query = ""

	def testScopeProperties(self):
		"""Test the scope properties from the .ini file"""

		self.assertEqual(self.view.scope_id, "scope")
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
