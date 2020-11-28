#!/usr/bin/env bash

# Metricator build script.

source ./scripts/shell_helpers/get_git_data.sh

WHATTOBUILD=$1

LINKERFLAGS="\
-X go.dev.pztrn.name/metricator/internal/common.Branch=${BRANCHNAME} \
-X go.dev.pztrn.name/metricator/internal/common.Build=${BUILDID} \
-X go.dev.pztrn.name/metricator/internal/common.CommitHash=${COMMITHASH} \
-X go.dev.pztrn.name/metricator/internal/common.Version=${VERSION}"


cd cmd/${WHATTOBUILD}
go build -tags netgo -ldflags "${LINKERFLAGS} -w -extldflags '-static'" -o ../../._bin/${WHATTOBUILD}
