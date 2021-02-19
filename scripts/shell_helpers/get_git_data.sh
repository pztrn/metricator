# Gets git data.
# Should be sourced where neccessary.
export BRANCHNAME=${BRANCHNAME:=$(git branch --no-color --show-current)}
export BUILDID=${BUILDID:=$(git rev-list HEAD --count)}
export COMMITHASH=${COMMITHASH:=$(git rev-parse --verify HEAD)}
export VERSION=${VERSION:=$(cat VERSION)}
