# build

Dockerfiles to build the deej app for different distros.

Top-level directory for each distro:
  - Dockerfile
  - docker-compose.yml
  - build.sh
  - install

Requirements to build:
  - docker
  - docker-compose

To build deej for a specific distro, run `./build.sh` in the distro's directory.

## debian
The built app should probably work on debian based distros with gnome desktop.
