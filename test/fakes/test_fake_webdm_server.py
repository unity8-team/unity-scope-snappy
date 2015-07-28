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

import testtools
import fixtures
import json
import requests

# Local imports
import test.fakes.fake_webdm_server as fake
from test.test_fixtures import ServerFixture

class TestFakeWebdmServer(testtools.TestCase, fixtures.TestWithFixtures):
	"""Test fake WebDM server to verify that it works as expected."""

	def setUp(self):
		"""Setup the test fixtures and run the store scope using the harness.
		"""

		super().setUp()

		server = ServerFixture(fake.FakeWebdmServer)
		self.useFixture(server)
		self.url = server.url[:-1]

	def testListStorePackages(self):
		"""Test the package list API."""

		response = requests.get(self.url + fake.PACKAGE_LIST_PATH)
		results = response.json()

		self.assertEqual(len(results), 2,
		                 "There should only be two packages")

		# Order is not enforced on the list, so we need to check it this
		# way.
		package1Found = False
		package2Found = False
		for package in results:
			if package["id"] == "package1.canonical":
				package1Found = True
			elif package["id"] == "package2.canonical":
				package2Found = True
			else:
				self.fail("Unexpected package: {}".format(package["id"]))

		if not package1Found:
			self.fail("Expected package with ID: \"package1.canonical\"")

		if not package2Found:
			self.fail("Expected package with ID: \"package2.canonical\"")

	def testListInstalledPackages(self):
		"""Test the package list API with a filter for installed packages."""

		response = requests.get(self.url + fake.PACKAGE_LIST_PATH,
		                        params={"installed_only": "true"})
		results = response.json()

		self.assertEqual(len(results), 1,
		                 "There should be a single result")

		self.assertEqual(results[0]["id"], "package1.canonical",
		                 "Expected the installed package to be package1")

	def testQueryInstalled(self):
		"""Test the query API with an installed package."""

		response = requests.get(self.url + fake.PACKAGE_LIST_PATH +
		                        "package1.canonical")
		result = response.json()
		self.assertEqual(result["id"], "package1.canonical")

	def testQueryNotInstalled(self):
		"""Test the query API with a package that is not installed."""

		response = requests.get(self.url + fake.PACKAGE_LIST_PATH +
		                        "package2.canonical")
		result = response.json()
		self.assertEqual(result["id"], "package2.canonical")

	def testQueryInstalledAndInstalledOnly(self):
		"""Test the query API with a filter for installed packages.

		Test with an installed package.
		"""

		response = requests.get(self.url + fake.PACKAGE_LIST_PATH +
		                        "package1.canonical",
		                         params={"installed_only": "true"})
		result = response.json()
		self.assertEqual(result["id"], "package1.canonical")

	def testQueryNotInstalledAndInstalledOnly(self):
		"""Test the query API with a filter for installed packages.

		Test with a package that is not installed.
		"""

		response = requests.get(self.url + fake.PACKAGE_LIST_PATH +
		                        "package2.canonical",
		                         params={"installed_only": "true"})
		self.assertEqual(response.status_code, 404,
		                 "Requesting \"package2\" while also limiting the search to installed packages should result in not found")

	def testPackageInstallation(self):
		"""Test the install API."""

		self.installPackage("package2.canonical")

	def testPackageUninstallation(self):
		"""Test the uninstall API."""

		self.uninstallPackage("package1.canonical")

	def testPackageInstallUninstallReinstall(self):
		"""Test a package install/uninstall/reinstall flow."""

		# Install package.
		self.installPackage("package2.canonical")

		# Now uninstall the package.
		self.uninstallPackage("package2.canonical")

		# Finally, reinstall the package.
		self.installPackage("package2.canonical")

	def installPackage(self, packageId):
		"""Install a package via the install API, making sure it succeeds.
		"""

		response = requests.put(self.url + fake.PACKAGE_LIST_PATH + packageId)
		self.assertEqual(response.status_code, 202)

		progress = 0
		status = ""
		while True:
			response = requests.get(self.url + fake.PACKAGE_LIST_PATH +
			                        packageId)
			result = response.json()
			if result["status"] == "installed":
				break

			newProgress = result["progress"]

			if newProgress <= progress:
				self.fail("Progress was {}, but last call it was {}. It should be increasing with each call".format(newProgress, progress))

			progress = newProgress

	def uninstallPackage(self, packageId):
		"""Uninstall a package via the uninstall API, making sure it succeeds.
		"""

		response = requests.delete(self.url + fake.PACKAGE_LIST_PATH +
		                        packageId)
		self.assertEqual(response.status_code, 202)

		progress = 0
		status = ""
		while True:
			response = requests.get(self.url + fake.PACKAGE_LIST_PATH +
			                        packageId)
			result = response.json()
			if result["status"] == "uninstalled":
				break

			newProgress = result["progress"]

			if newProgress <= progress:
				self.fail("Progress was {}, but last call it was {}. It should be increasing with each call".format(newProgress, progress))

			progress = newProgress
