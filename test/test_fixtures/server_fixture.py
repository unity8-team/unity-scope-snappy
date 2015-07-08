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
