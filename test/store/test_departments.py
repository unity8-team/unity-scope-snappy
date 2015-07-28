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
