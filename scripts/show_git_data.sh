#!/usr/bin/env bash

# Showing git data in console.
if [ -d .git ]; then
    source ./scripts/shell_helpers/get_git_data.sh
else
    source ./scripts/shell_helpers/get_release_data.sh
fi

echo "* Branch: ${BRANCHNAME}"
echo "* Build ID: ${BUILDID}"
echo "* Commit hash: ${COMMITHASH}"
echo "* Version: ${VERSION}"
