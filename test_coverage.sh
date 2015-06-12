#!/bin/bash

tmp_directory="$(mktemp -d)"
coverage_final=$tmp_directory/coverage.out
coverage_scope=$tmp_directory/coverage_scope.out
coverage_webdm=$tmp_directory/coverage_webdm.out
coverage_daemon=$tmp_directory/coverage_daemon.out

# Setup final coverage file
echo "mode: set" > $coverage_final

# Add the scope's coverage to the overarching coverage file
go test launchpad.net/unity-scope-snappy/scope -coverprofile=$coverage_scope
cat $coverage_scope | grep -v "mode: set" >> $coverage_final

# Add webdm's coverage to the overarching coverage file
go test launchpad.net/unity-scope-snappy/webdm -coverprofile=$coverage_webdm
cat $coverage_webdm | grep -v "mode: set" >> $coverage_final

# Add webdm's coverage to the overarching coverage file
go test launchpad.net/unity-scope-snappy/progress-daemon/daemon -coverprofile=$coverage_daemon
cat $coverage_daemon | grep -v "mode: set" >> $coverage_final

if [ "$1" == "xml" ]; then
	gocov convert $coverage_final | gocov-xml > coverage.xml
elif [ "$1" == "html" ]; then
	go tool cover -html=$coverage_final
fi
