import http.server
import urllib
import json

class FakeWebdmServerHandler(http.server.BaseHTTPRequestHandler):
	_PACKAGE_LIST_PATH = "/api/v2/packages/"

	_PACKAGE1 = {
		"id": "package1.canonical",
		"name": "package1",
		"vendor": "Canonical",
		"version": "0.1",
		"description": "description1",
		"icon": "http://icon1",
		"type": "app",
		"status": "installed",
		"installed_size": 123456
	}

	_PACKAGE2 = {
		"id": "package2.canonical",
		"name": "package2",
		"vendor": "Canonical",
		"version": "0.2",
		"description": "description2",
		"icon": "http://icon2",
		"type": "app",
		"status": "uninstalled",
		"download_size": 123456
	}

	def sendJson(self, code, json):
		self.send_response(code)
		self.send_header('Content-Type', 'application/json')
		self.end_headers()
		self.wfile.write(json.encode())

	def do_GET(self):
		parsedPath = urllib.parse.urlparse(self.path)
		query = urllib.parse.parse_qs(parsedPath.query)
		installedOnly = False

		if "installed_only" in query and query["installed_only"][0] == "true":
			installedOnly = True

		if parsedPath.path.startswith(self._PACKAGE_LIST_PATH):
			packageId = parsedPath.path[len(self._PACKAGE_LIST_PATH):]

			# If no package ID was provided, list packages instead
			if len(packageId) == 0:
				if installedOnly:
					self.sendJson(200, json.dumps([self._PACKAGE1]))
				else:
					self.sendJson(200, json.dumps([self._PACKAGE1, self._PACKAGE2]))
			else:
				if packageId == "package1.canonical":
					self.sendJson(200, json.dumps(self._PACKAGE1))
				elif packageId == "package2.canonical" and not installedOnly:
					self.sendJson(200, json.dumps(self._PACKAGE2))
				else:
					self.send_error(404, "snappy package not found {}\n".format(packageId))

			return

		raise NotImplementedError(self.path)

class FakeWebdmServer(http.server.HTTPServer):
	def __init__(self, address):
		super().__init__(address, FakeWebdmServerHandler)
