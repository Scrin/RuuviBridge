package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type GatewayPolling struct {
	Enabled    *bool         `yaml:"enabled,omitempty"`
	GatewayUrl string        `yaml:"gateway_url"`
	Interval   time.Duration `yaml:"interval"`
}

type MQTTListener struct {
	Enabled       *bool  `yaml:"enabled,omitempty"`
	BrokerAddress string `yaml:"broker_address"`
	BrokerPort    int    `yaml:"broker_port"`
	ClientID      string `yaml:"client_id"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	TopicPrefix   string `yaml:"topic_prefix"`
}

type Processing struct {
	ExtendedValues *bool    `yaml:"extended_values,omitempty"`
	FilterMode     string   `yaml:"filter_mode"`
	FilterList     []string `yaml:"filter_list"`
}

type InfluxDBPublisher struct {
	Enabled         *bool         `yaml:"enabled,omitempty"`
	MinimumInterval time.Duration `yaml:"minimum_interval,omitempty"`
	Url             string        `yaml:"url"`
	AuthToken       string        `yaml:"auth_token"`
	Org             string        `yaml:"org"`
	Bucket          string        `yaml:"bucket"`
	Measurement     string        `yaml:"measurement"`
}

type Prometheus struct {
	Enabled *bool `yaml:"enabled,omitempty"`
	Port    int   `yaml:"port"`
}

type MQTTPublisher struct {
	Enabled                      *bool         `yaml:"enabled,omitempty"`
	MinimumInterval              time.Duration `yaml:"minimum_interval,omitempty"`
	BrokerAddress                string        `yaml:"broker_address"`
	BrokerPort                   int           `yaml:"broker_port"`
	ClientID                     string        `yaml:"client_id"`
	Username                     string        `yaml:"username"`
	Password                     string        `yaml:"password"`
	TopicPrefix                  string        `yaml:"topic_prefix"`
	HomeassistantDiscoveryPrefix string        `yaml:"homeassistant_discovery_prefix,omitempty"`
}

type Logging struct {
	Type       string `yaml:"type"`
	Level      string `yaml:"level"`
	Timestamps *bool  `yaml:"timestamps,omitempty"`
	WithCaller bool   `yaml:"with_caller,omitempty"`
}

type Config struct {
	GatewayPolling    *GatewayPolling    `yaml:"gateway_polling,omitempty"`
	MQTTListener      *MQTTListener      `yaml:"mqtt_listener,omitempty"`
	Processing        *Processing        `yaml:"processing,omitempty"`
	InfluxDBPublisher *InfluxDBPublisher `yaml:"influxdb_publisher,omitempty"`
	Prometheus        *Prometheus        `yaml:"prometheus,omitempty"`
	MQTTPublisher     *MQTTPublisher     `yaml:"mqtt_publisher,omitempty"`
	TagNames          map[string]string  `yaml:"tag_names,omitempty"`
	Logging           Logging            `yaml:"logging"`
	Debug             bool               `yaml:"debug"`
}

func ReadConfig(configFile string, strict bool) (Config, error) {
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return Config{}, errors.New(fmt.Sprintf("No config found! Tried to open \"%s\"", configFile))
	}

	f, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var conf Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(strict)
	err = decoder.Decode(&conf)

	if err != nil {
		return Config{}, err
	}
	return conf, nil
}
