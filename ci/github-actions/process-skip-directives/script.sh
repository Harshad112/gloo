#!/bin/bash

################################################################
# This script checks to see if a changelog file has been added
# with the "skipCI" field set to true.
#
# It will set the SKIP_CI_TESTS env variable to true which can be
# used across steps in the same job
#
# It will  write a file called skip-ci.txt with contents of
# either SKIP_CI_TESTS=true or SKIP_CI_TESTS=false - which can be used
# across different jobs in the same workflow
################################################################

set -ex

echo "SKIP_CI_TESTS=false" > skip-ci.txt

if [[ $(git diff origin/main HEAD --name-only | grep "changelog/" | wc -l) = "1" ]]; then
    echo "exactly one changelog added since main"
    changelogFileName=$(git diff origin/main HEAD --name-only | grep "changelog/")
    echo "changelog file name == $changelogFileName"
    if [[ $(cat $changelogFileName | grep "skipCI: true") ]]; then
        echo "skip CI"
	echo "SKIP_CI_TESTS=true" > skip-ci.txt
    else
        echo "skipCI not set to true - not skipping CI"
    fi
else
    echo "no changelog found (or more than one changelog found) - not skipping CI"
fi