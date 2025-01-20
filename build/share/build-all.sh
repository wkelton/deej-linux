#!/bin/bash

set -euo

#git config --add safe.directory .

export GOPATH=$(realpath ./go)
export GOCACHE=$(realpath ./.cache)
mkdir "${GOCACHE}"

pkg/deej/scripts/linux/build-dev.sh
pkg/deej/scripts/linux/build-release.sh
