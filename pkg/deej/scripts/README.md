## Developer scripts

This document lists the various scripts in the project and their purposes.

> The Windows scripts are intended to be run from the repository root. The Linux scripts can be run from any working directory.

### Windows

- [`build-dev.bat`](./windows/build-dev.bat): Builds deej with a console window, for development purposes
- [`build-release.bat`](./windows/build-release.bat): Builds deej as a standalone tray application without a console window, for releases
- [`build-all.bat`](./windows/build-all.bat): Helper script to build all variants
- [`make-icon.bat`](./windows/make-icon.bat): Converts a .ico file to an icon byte array in a Go file. Used by our systray library. You shouldn't need to run this unless you change the deej logo
- [`make-rsrc.bat`](./windows/make-rsrc.bat): Generates a `rsrc.syso` resource file inside `cmd` alongside `main.go` - This indicates to the Go linker to use the deej application manifest and icon when building.
- [`prepare-release.bat`](./windows/prepare-release.bat): Tags, builds and renames the release binaries in preparation for a GitHub release. Usage: `prepare-release.bat vX.Y.Z` (binaries will be under `releases\vX.Y.Z\`)

### Linux

- [`build-dev.sh`](./linux/build-dev.sh): Builds deej for development purposes
- [`build-release.sh`](./linux/build-release.sh): Builds deej for releases
- [`build-all.sh`](./linux/build-all.sh): Helper script to build all variants
- [`build-cross-windows.sh`](./linux/build-cross-windows.sh): Cross-compiles the development build for Windows/amd64
