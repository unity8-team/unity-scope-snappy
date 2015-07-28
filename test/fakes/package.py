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
