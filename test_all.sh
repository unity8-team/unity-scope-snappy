#!/bin/bash

directory=$(dirname $0)

# Run golang tests
"$directory/test_coverage.sh" $@
returnValue=$?
if [ $returnValue -ne 0 ]; then
	exit $returnValue # Tests failed, so this script should also fail
fi

# Run Python integration tests
"$directory/test_integration.sh"
returnValue=$?
if [ $returnValue -ne 0 ]; then
	exit $returnValue # Tests failed, so this script should also fail
fi
