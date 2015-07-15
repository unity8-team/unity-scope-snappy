# Go parameters
GOCMD=go
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GOFMT=gofmt -w

TOP_PACKAGE := launchpad.net/unity-scope-snappy
EXECUTABLES := store progress-daemon
PACKAGES_TO_TEST := progress-daemon/daemon \
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
test: go_tests integration_tests
go_tests: $(TEST_LIST)
fmt: $(FMT_LIST)

.PHONY: coverage
coverage:
	./test_coverage.sh $(ABSOLUTE_PACKAGES_TO_TEST)

.PHONY: coverage_html
coverage_html:
	./test_coverage.sh -t html $(ABSOLUTE_PACKAGES_TO_TEST)

.PHONY: coverage_xml
coverage_xml:
	./test_coverage.sh -t xml $(ABSOLUTE_PACKAGES_TO_TEST)

.PHONY: integration_tests
integration_tests:
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
