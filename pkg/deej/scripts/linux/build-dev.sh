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
BUILD_TYPE=dev
echo 'Embedding build-time parameters:'
echo "- gitCommit $GIT_COMMIT"
echo "- versionTag $VERSION_TAG"
echo "- buildType $BUILD_TYPE"

go -C "$REPO_ROOT" build -o install/deej-dev -ldflags "-X main.gitCommit=$GIT_COMMIT -X main.versionTag=$VERSION_TAG -X main.buildType=$BUILD_TYPE" ./pkg/deej/cmd
if [ $? -eq 0 ]; then
    echo 'Done.'
else
    echo 'Error: "go build" exited with a non-zero code.'
    exit 1
fi
