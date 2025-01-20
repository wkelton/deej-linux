# deej-linux

This is a fork of the deej project to maintain tweaks required for easy use on Linux. See [Original README](./docs/README.md) / [Upstream README](https://github.com/omriharel/deej/blob/master/README.md).

## Changes
- Tweak for apps whose audio goes through pipewire
- Updated go packages
- Support for $XDG_CONFIG_HOME/deej/config.yaml
- Configurable adc_size
- config.yaml with linux-sensible defaults
- Build scripts place executables in `install` directory

## Requirements
You'll need the permissions to access the serial device. Something like:
```
sudo usermod -aG dialout <username> # requires re-login
```

## Building

See [build/README.md](./build/README.md)
