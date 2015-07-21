class Package:
	def __init__(self, id="", name="", vendor="", version="", description="", icon="", type="", status="", installed_size=0, download_size=0, progress=0):
		self.id = id
		self.name = name
		self.vendor = vendor
		self.version = version
		self.description = description
		self.icon = icon
		self.type = type
		self.status = status
		self.installed_size = installed_size
		self.download_size = download_size
		self.progress = progress
