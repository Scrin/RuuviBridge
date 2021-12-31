# RuuviBridge

RuuviBridge is designed to act as a "data bridge" between various sources and consumers of data. Original design goal is to work as a drop-in replacement for [RuuviCollector](https://github.com/Scrin/RuuviCollector) for users who have a [Ruuvi Gateway](https://ruuvi.com/gateway/) or use [ruuvi-go-gateway](https://github.com/Scrin/ruuvi-go-gateway).

### Features

Supports following sources (sources of RuuviTag data):

- MQTT (in Ruuvi Gateway format)
- Ruuvi Gateway by polling the /history http-api endpoint
- HTTP POST (in Ruuvi Gateway format, the custom http server setting)

Supports following sinks (things that use the data):

- InfluxDB 1.8 and 2.x
- Prometheus
- MQTT (including Home Assistant MQTT discovery for automatic configuration)

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

### Configuration

Check [config.sample.yml](./config.sample.yml) for a sample config. By default the bridge assumes to find a file called `config.yml` in the current working directory, but that can be overridden with `-config /path/to/config.yml` command line flag.

By default RuuviBridge parses the config in a flexible way, ignoring all unknown fields. This can be changed with `-strict-config` command line flag, which will make RuuviBridge throw errors if there are unknown entries in the config. Do note that this only validates whether the config has a valid structure with right keys (ie. no typos in the keys), it does not validate whether the config makes sense as such.

### Installation

Recommended method is using Docker with the prebuilt dockerimage: [ghcr.io/scrin/ruuvibridge](https://ghcr.io/scrin/ruuvibridge) for which you can use the provided [composefile](./docker-compose.yml)

Without docker you can download prebuilt binaries from the [releases](https://github.com/Scrin/RuuviBridge/releases) page. For production use it's recommended to set up as a service.

### Home Assistant MQTT discovery

Home Assistant allows automatic configuration of MQTT entities using [MQTT Discovery](https://www.home-assistant.io/docs/mqtt/discovery/). To enable RuuviBridge to automatically configure all of your RuuviTags to Home Assistant for you, all you need to do (assuming default configuration) is to set `homeassistant_discovery_prefix` in the config under `mqtt_publisher`. In default Home Assistant configuration this should be simply `homeassistant`.

After setting this configuration, it should be a matter of seconds before your RuuviTags should appear as devices in Home Assistant for reporting all available measurements, with properly set names, units, icons and other attributes.
