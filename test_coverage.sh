#!/bin/bash

coverage_final="$(mktemp)"

# get_coverage runs the given tests and adds the results to the final coverage
# file.
#
# Parameters:
# $1: Go path under test.
# $2: Final coverage report.
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

# Setup final coverage file
echo "mode: set" > $coverage_final

# Add the store scope's coverage to the final coverage file
get_coverage "launchpad.net/unity-scope-snappy/store" $coverage_final

# Add the package manager's coverage to the final coverage file
get_coverage "launchpad.net/unity-scope-snappy/store/packages" $coverage_final

# Add webdm's coverage to the final coverage file
get_coverage "launchpad.net/unity-scope-snappy/webdm" $coverage_final

# Add progress daemon's coverage to the final coverage file
get_coverage "launchpad.net/unity-scope-snappy/progress-daemon/daemon" $coverage_final

if [ "$1" == "xml" ]; then
	gocov convert $coverage_final | gocov-xml > coverage.xml
elif [ "$1" == "html" ]; then
	go tool cover -html=$coverage_final
fi
