# deej-linux

This is a fork of the deej project to maintain tweaks required for easy use on Linux. See [Original README](./docs/README.md) / [Upstream README](https://github.com/omriharel/deej/blob/master/README.md).

## Changes

### Fixes
- Changed pulseaudio session finder to match application on several fields
  - Should help with applications run with pipewire or wrapped in proton/wine
- Merged reconnect fix from upstream (mainly for Windows)

### Config
- Support for $XDG_CONFIG_HOME/deej/config.yaml
- config.yaml with linux-sensible defaults
- Configurable adc_size
- Settings for pulseaudio application matching

### Building
- Build scripts place executables in `install` directory
- Build with docker containers - Debian 13 currently
- Addded script to cross build for Windows from Linux

### Dependencies
- Updated go packages
- Updated go version with required code changes

## Requirements
You'll need the permissions to access the serial device. Something like:
```
sudo usermod -aG dialout <username> # requires re-login
```

Currently built for go 1.26.

## Custom config settings

```yaml
# setting for the adc size of your device; only need to change if your device isn't sending values [0,1023]
adc_size: 1024

# custom settings for linux
linux:
  # look at additional pulseaudio fields when matching applications with names in slider_mapping
  use_extra_app_name_fields: true
  # log the names found for each pulseaudio application
  log_session_names: true
  # log all the pulseaudio properties for each pulseaudio application
  log_session_properties: false
```

## Building

See [build/README.md](./build/README.md)
