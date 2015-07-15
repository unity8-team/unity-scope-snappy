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

coverage_type=""

while getopts ":t:" opt; do
	case $opt in
		t)
			coverage_type=$OPTARG
			;;
		\?)
			echo "Invalid option: -$OPTARG" >&2
			;;
		:)
			echo "Option -$OPTARG requires an argument." >&2
			exit 1
			;;
	esac
done

shift $((OPTIND-1))

# Setup final coverage file
echo "mode: set" > $coverage_final

get_coverages $coverage_final $@

echo "Cov type: $coverage_type"

if [ "$coverage_type" == "xml" ]; then
	gocov convert $coverage_final | gocov-xml > coverage.xml
elif [ "$coverage_type" == "html" ]; then
	go tool cover -html=$coverage_final
fi
