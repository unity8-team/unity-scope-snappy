#!/bin/bash

coverage_final="$(mktemp)"

# get_coverage runs the given tests and adds the results to the final coverage
# file.
#
# Parameters:
# $1: Go path under test.
# $2: Final coverage file.
function get_coverage
{
	local go_path=$1
	local coverage_final=$2
	local coverage_tmp="$(mktemp)"

	go test $go_path -coverprofile=$coverage_tmp
	local returnValue=$?
	if [ $returnValue -eq 0 ]; then
		# Tests passed, so add this to the final coverage file
		cat $coverage_tmp | grep -v "mode: set" >> $coverage_final
	else
		exit $returnValue # Tests failed, so this script should also fail
	fi
}

# get_coverages accepts an array of Go packages and gets coverage for each,
# placing the results in the requested file.
#
# Parameters:
# $1: Final coverage file.
# $2-$n: Set of Go packages.
function get_coverages
{
	local coverage_final=$1

	for go_package in "${@:2}"; do
		get_coverage $go_package $coverage_final
	done
}

# Setup final coverage file
echo "mode: set" > $coverage_final

get_coverages $coverage_final \
	"launchpad.net/unity-scope-snappy/progress-daemon/daemon" \
	"launchpad.net/unity-scope-snappy/store/actions" \
	"launchpad.net/unity-scope-snappy/store/packages" \
	"launchpad.net/unity-scope-snappy/store/previews" \
	"launchpad.net/unity-scope-snappy/store/previews/humanize" \
	"launchpad.net/unity-scope-snappy/store/previews/interfaces" \
	"launchpad.net/unity-scope-snappy/store/previews/packages" \
	"launchpad.net/unity-scope-snappy/store/previews/packages/templates" \
	"launchpad.net/unity-scope-snappy/store/utilities" \
	"launchpad.net/unity-scope-snappy/webdm"

if [ "$1" == "xml" ]; then
	gocov convert $coverage_final | gocov-xml > coverage.xml
elif [ "$1" == "html" ]; then
	go tool cover -html=$coverage_final
fi
