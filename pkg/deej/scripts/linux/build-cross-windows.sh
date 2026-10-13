#!/bin/sh

SCRIPT_PATH=$0
case "$SCRIPT_PATH" in
    /*) ;;
    */*) SCRIPT_PATH=$PWD/$SCRIPT_PATH ;;
    *) SCRIPT_PATH=$(command -v "$SCRIPT_PATH") || exit 1 ;;
esac
SCRIPT_DIR=$(dirname -- "$SCRIPT_PATH") || exit 1
REPO_ROOT=$SCRIPT_DIR/../../../../

echo 'Building deej (development)...'

# shove git commit, version tag into env
GIT_COMMIT=$(git -C "$REPO_ROOT" rev-list -1 --abbrev-commit HEAD)
VERSION_TAG=$(git -C "$REPO_ROOT" describe --tags --always)
echo 'Embedding build-time parameters:'
echo "- gitCommit $GIT_COMMIT"
echo "- versionTag $VERSION_TAG"

export GOOS=windows
export GOARCH=amd64

if ! go -C "$REPO_ROOT" build -o install/deej-dev.exe -ldflags "-X main.gitCommit=$GIT_COMMIT -X main.versionTag=$VERSION_TAG -X main.buildType=dev" ./pkg/deej/cmd; then
    echo 'Error: failed to build the Windows development executable.'
    exit 1
fi

if ! go -C "$REPO_ROOT" build -o install/deej-release.exe -ldflags "-H=windowsgui -s -w -X main.gitCommit=$GIT_COMMIT -X main.versionTag=$VERSION_TAG -X main.buildType=release" ./pkg/deej/cmd; then
    echo 'Error: failed to build the Windows release executable.'
    exit 1
fi

echo 'Done.'
