# RuuviBridge

RuuviBridge is designed to act as a "data bridge" between various sources and consumers of data. Original design goal is to work as a drop-in replacement for [RuuviCollector](https://github.com/Scrin/RuuviCollector) for users who have a [Ruuvi Gateway](https://ruuvi.com/gateway/).

Note: This is very early in development; breaking changes will occur. Be sure to check the config.sample.yml for changes when you update

### Features

Supports following sources (sources of RuuviTag data):

- MQTT (in RuuviGateway format)
- RuuviGateway /history http-api endpoint

Supports following sinks (things that use the data):

- InfluxDB 1.8 and 2.x
- Prometheus

Supports following RuuviTag [Data Formats](https://github.com/ruuvi/ruuvi-sensor-protocols):

- Data Format 3: "RAW v1" BLE Manufacturer specific data, all current sensor readings
- Data Format 5: "RAW v2" BLE Manufacturer specific data, all current sensor readings + extra

Supports following data from the tag (depending on tag firmware):

- Temperature (Celsius)
- Relative humidity (0-100%)
- Air pressure (Pascal)
- Acceleration for X, Y and Z axes (g)
- Battery voltage (Volts)
- TX power (dBm)
- RSSI (Signal strength _at the receiver_, dBm)
- Movement counter (Running counter incremented each time a motion detection interrupt is received)
- Measurement sequence number (Running counter incremented each time a new measurement is taken on the tag)

Ability to calculate following values in addition to the raw data (the accuracy of these values are approximations):

- Total acceleration (g)
- Absolute humidity (g/m³)
- Dew point (Celsius)
- Equilibrium vapor pressure (Pascal)
- Air density (Accounts for humidity in the air, kg/m³)
- Acceleration angle from X, Y and Z axes (Degrees)

### Roadmap

In no particular order:

- Proper documentation
- Standalone binary releases
- Properly versioned releases with changelogs
- Support for MQTT as an output for making it easier to use your own applications with already parsed data
- HTTP endpoint to allow "pushes" from a Ruuvi Gateway without having a MQTT server
- And other stuff I forgot I had plans for

### Configuration

Check [config.sample.yml](./config.sample.yml) for a sample config. By default the bridge assumes to find a file called `config.yml` in the current working directory, but that can be overridden with `-config /path/to/config.yml` command line flag.

### Installation

Recommended method is using Docker with the prebuilt dockerimage: [ghcr.io/scrin/ruuvibridge](https://ghcr.io/scrin/ruuvibridge) for which you can use the provided [composefile](./docker-compose.yml)

Without docker you have to build the binaries yourself until I have time to set up a release process. Easiest way to do this is to install Go 1.17 or later and run `go install -v github.com/Scrin/RuuviBridge/cmd/ruuvibridge@latest` which will download, build and install ruuvibridge into `$GOPATH/bin`
