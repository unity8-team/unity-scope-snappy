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

import fixtures
import multiprocessing

class ServerFixture(fixtures.Fixture):
	def __init__(self, serverClass, *serverArguments):
		super().__init__()
		self.serverClass = serverClass
		self.serverArguments = serverArguments

	def setUp(self):
		super().setUp()
		self._startServer()

	def _startServer(self):
		queue = multiprocessing.Queue()
		process = multiprocessing.Process(target=self._serverProcess,
		                                  args=(queue, self.serverClass))
		process.start()
		self.addCleanup(self._stopServer, process)
		self.url = queue.get()

	def _stopServer(self, process):
		process.terminate()
		process.join()

	def _serverProcess(self, queue, serverClass):
		address = ('localhost', 0)
		server = serverClass(address, *self.serverArguments)
		server.url = 'http://localhost:{}/'.format(server.server_port)
		queue.put(server.url)
		server.serve_forever()
