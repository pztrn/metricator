#!/usr/bin/env bash

# Metricator build script.

GO=${GO:=$(which go)}

if [ -d .git ]; then
    source ./scripts/shell_helpers/get_git_data.sh
else
    source ./scripts/shell_helpers/get_release_data.sh
fi

WHATTOBUILD=$1

LINKERFLAGS="\
-X go.dev.pztrn.name/metricator/internal/common.Branch=${BRANCHNAME} \
-X go.dev.pztrn.name/metricator/internal/common.Build=${BUILDID} \
-X go.dev.pztrn.name/metricator/internal/common.CommitHash=${COMMITHASH} \
-X go.dev.pztrn.name/metricator/internal/common.Version=${VERSION}"


echo "Using $(go version) at ${GO}"

cd cmd/${WHATTOBUILD}
${GO} build -tags netgo -ldflags "${LINKERFLAGS} -w -extldflags '-static'" -o ../../._bin/${WHATTOBUILD}
