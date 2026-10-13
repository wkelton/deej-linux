#!/bin/sh

SCRIPT_PATH=$0
case "$SCRIPT_PATH" in
    /*) ;;
    */*) SCRIPT_PATH=$PWD/$SCRIPT_PATH ;;
    *) SCRIPT_PATH=$(command -v "$SCRIPT_PATH") || exit 1 ;;
esac
SCRIPT_DIR=$(dirname -- "$SCRIPT_PATH") || exit 1

echo 'Building deej (all)...'

"$SCRIPT_DIR/build-dev.sh" || exit $?
"$SCRIPT_DIR/build-release.sh" || exit $?
