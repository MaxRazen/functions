#!/usr/bin/env bash
# The idea was taken from
# https://github.com/google/go-cloud/blob/master/internal/testing/runchecks.sh

echo
echo "************************"
echo "* Running tests for all the modules"
echo "************************"
echo
while read -r path || [[ -n "$path" ]]; do
    echo "Module: $path"
    ( cd "$path" && go test -count=1 && echo "  OK" ) || { echo "FAIL: tests run failed"; }
done < <( sed -e '/^#/d' -e '/^$/d' allmodules | awk '{print $1}' )
# The above filters out comments and empty lines from allmodules and only takes
# the first (whitespace-separated) field from each line.