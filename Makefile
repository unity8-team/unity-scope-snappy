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

# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GOFMT=gofmt -w

TOP_PACKAGE := launchpad.net/unity-scope-snappy
EXECUTABLES := store package-management-daemon
PACKAGES_TO_TEST := package-management-daemon/daemon \
                    store/actions \
                    store/packages \
                    store/packages/fakes \
                    store/previews \
                    store/previews/humanize \
                    store/previews/fakes \
                    store/previews/packages \
                    store/previews/packages/templates \
                    store/utilities \
                    webdm

ALL_LIST = $(EXECUTABLES) $(PACKAGES_TO_TEST)

ABSOLUTE_PACKAGES_TO_TEST = $(foreach item, $(PACKAGES_TO_TEST), $(TOP_PACKAGE)/$(item))

INSTALL_LIST = $(foreach item, $(EXECUTABLES), $(item)_install)
TEST_LIST = $(foreach item, $(PACKAGES_TO_TEST), $(item)_test)
FMT_LIST = $(foreach item, $(ALL_LIST), $(item)_fmt)

.PHONY: $(INSTALL_LIST) $(TEST_LIST) $(FMT_LIST)

all: install
install: $(INSTALL_LIST)
check: test
test: go_tests integration_tests
go_tests: $(TEST_LIST)
fmt: $(FMT_LIST)

.PHONY: coverage
coverage:
	./test_coverage.sh $(ABSOLUTE_PACKAGES_TO_TEST)

.PHONY: coverage.html
coverage.html:
	./test_coverage.sh -t html $(ABSOLUTE_PACKAGES_TO_TEST)

.PHONY: coverage.xml
coverage.xml:
	./test_coverage.sh -t xml $(ABSOLUTE_PACKAGES_TO_TEST)

.PHONY: integration_tests
integration_tests: install
	python3 -m unittest discover -s test

$(BUILD_LIST): %_build: %_fmt
	$(GOBUILD) $(TOP_PACKAGE)/$*
$(CLEAN_LIST): %_clean:
	$(GOCLEAN) $(TOP_PACKAGE)/$*
$(INSTALL_LIST): %_install:
	$(GOINSTALL) $(TOP_PACKAGE)/$*
$(TEST_LIST): %_test:
	$(GOTEST) $(TOP_PACKAGE)/$*
$(FMT_LIST): %_fmt:
	$(GOFMT) ./$*
