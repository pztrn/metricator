#!/usr/bin/env bash

# Showing git data in console.
source ./scripts/shell_helpers/get_git_data.sh

echo "* Branch: ${BRANCHNAME}"
echo "* Build ID: ${BUILDID}"
echo "* Commit hash: ${COMMITHASH}"
echo "* Version: ${VERSION}"
