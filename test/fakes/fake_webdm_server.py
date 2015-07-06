import http.server
import multiprocessing
import urllib
import json

# Local imports
from .package import Package

PACKAGE_LIST_PATH = "/api/v2/packages/"

def finishOperation(status):
	if status == "installing":
		return "installed"
	elif status == "uninstalling":
		return "uninstalled"

	return status

def undoOperation(status):
	if status == "installing":
		return "uninstalled"
	elif status == "uninstalling":
		return "installed"

	return status

class FakeWebdmServerHandler(http.server.BaseHTTPRequestHandler):
	_PROGRESS_STEP = 50

	def sendJson(self, code, json):
		self.send_response(code)
		self.send_header('Content-Type', 'application/json')
		self.end_headers()
		self.wfile.write(json.encode())

	def sendPackages(self, packages):
		jsonArray = []
		for package in packages:
			jsonArray.append(package.__dict__)

		self.sendJson(200, json.dumps(jsonArray))

	def findPackage(self, packageId):
		for index, package in enumerate(self.server.PACKAGES):
			if package.id == packageId:
				return (package, index)

		return (None, None)

	def continueOperation(self, package):
		if package.progress >= 100:
			if package.status == "installing":
				package.installed_size = package.download_size
				package.download_size = 0
			elif package.status == "uninstalling":
				package.download_size = package.installed_size
				package.installed_size = 0

			package.status = finishOperation(package.status)
			package.progress = 0
		else:
			package.progress += self._PROGRESS_STEP

		return package

	def do_GET(self):
		parsedPath = urllib.parse.urlparse(self.path)
		query = urllib.parse.parse_qs(parsedPath.query)
		installedOnly = False

		if "installed_only" in query and query["installed_only"][0] == "true":
			installedOnly = True

		if parsedPath.path.startswith(PACKAGE_LIST_PATH):
			packageId = parsedPath.path[len(PACKAGE_LIST_PATH):]

			# If no package ID was provided, list packages instead
			if len(packageId) == 0:
				packages = self.server.PACKAGES

				if installedOnly:
					packages = []
					for package in self.server.PACKAGES:
						if package.status == "installed":
							packages.append(package)

				self.sendPackages(packages)
			else:
				package, index = self.findPackage(packageId)
				if package != None:
					if (package.status == "installing" or
					    package.status == "uninstalling"):
						package = self.continueOperation(package)
						self.server.PACKAGES[index] = package

					if not installedOnly or (package.status == "installed"):
							self.sendJson(200, json.dumps(package.__dict__))
							return

				self.send_error(404, "snappy package not found {}\n".format(packageId))

			return

		raise NotImplementedError(self.path)

	def do_PUT(self):
		parsedPath = urllib.parse.urlparse(self.path)
		query = urllib.parse.parse_qs(parsedPath.query)

		if parsedPath.path.startswith(PACKAGE_LIST_PATH):
			packageId = parsedPath.path[len(PACKAGE_LIST_PATH):]

			if len(packageId) > 0:
				package, index = self.findPackage(packageId)

				if package != None and not self.server.ignoreRequests:
					package.status = "installing"
					self.server.PACKAGES[index] = package

				self.sendJson(202, json.dumps("Accepted"))
			else:
				self.send_error(500, "PUT here makes no sense")

			return

		raise NotImplementedError(self.path)

	def do_DELETE(self):
		parsedPath = urllib.parse.urlparse(self.path)
		query = urllib.parse.parse_qs(parsedPath.query)

		if parsedPath.path.startswith(PACKAGE_LIST_PATH):
			packageId = parsedPath.path[len(PACKAGE_LIST_PATH):]

			if len(packageId) > 0:
				package, index = self.findPackage(packageId)
				if package != None and not self.server.ignoreRequests:
					package.status = "uninstalling"
					self.server.PACKAGES[index] = package

				self.sendJson(202, json.dumps("Accepted"))
			else:
				self.send_error(500, "DELETE here makes no sense")

			return

		raise NotImplementedError(self.path)

class FakeWebdmServer(http.server.HTTPServer):
	def __init__(self, address, ignoreRequests=False):
		self.ignoreRequests = ignoreRequests

		manager = multiprocessing.Manager()
		self.PACKAGES = manager.list()
		self.PACKAGES.append(Package(
				id = "package1.canonical",
				name = "package1",
				vendor = "Canonical",
				version = "0.1",
				description = "description1",
				icon = "http://icon1",
				type = "app",
				status = "installed",
				installed_size = 123456
		))
		self.PACKAGES.append(Package(
				id = "package2.canonical",
				name = "package2",
				vendor = "Canonical",
				version = "0.2",
				description = "description2",
				icon = "http://icon2",
				type = "app",
				status = "uninstalled",
				download_size = 123456
		))

		super().__init__(address, FakeWebdmServerHandler)
