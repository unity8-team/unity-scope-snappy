#!/usr/bin/env python3

import os
import sys
import unittest
import fixtures

from scope_harness.testing import ScopeHarnessTestCase
from scope_harness import (ScopeHarness,
                           Parameters,
                           DepartmentMatcher,
                           ChildDepartmentMatcher)

# Local imports
from test.fakes import FakeWebdmServer
from test.test_fixtures import ServerFixture

class TestDepartments(ScopeHarnessTestCase, fixtures.TestWithFixtures):
	"""Test the various departments that are available in the store scope.
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

	def testStoreDepartment(self):
		"""Test department for in-store packages.

		Make sure that is has a child for installed packages.
		"""

		self.assertMatchResult(DepartmentMatcher()
			.has_exactly(1) # One child department
			.label("All Categories")
			.is_root(True)
			.is_hidden(False)
			.child(ChildDepartmentMatcher("installed")
				.label("My Snaps")
				.has_children(False)
				.is_active(False))
			.match(self.view.browse_department("")))

	def testInstalledDepartment(self):
		"""Test that the department for installed packages has no children.
		"""

		self.assertMatchResult(DepartmentMatcher()
			.has_exactly(0) # No child departments
			.label("My Snaps")
			.is_root(False)
			.is_hidden(False)
			.match(self.view.browse_department("installed")))
