package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type GatewayPolling struct {
	GatewayUrl string        `yaml:"gateway_url"`
	Interval   time.Duration `yaml:"interval"`
}

type MQTTListener struct {
	BrokerAddress string `yaml:"broker_address"`
	BrokerPort    int    `yaml:"broker_port"`
	ClientID      string `yaml:"client_id"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	TopicPrefix   string `yaml:"topic_prefix"`
}

type Processing struct {
	ExtendedValues *bool `yaml:"extended_values,omitempty"`
}

type InfluxDBPublisher struct {
	Url         string `yaml:"url"`
	AuthToken   string `yaml:"auth_token"`
	Org         string `yaml:"org"`
	Bucket      string `yaml:"bucket"`
	Measurement string `yaml:"measurement"`
}

type Config struct {
	Debug             bool               `yaml:"debug,omitempty"`
	GatewayPolling    *GatewayPolling    `yaml:"gateway_polling,omitempty"`
	MQTTListener      *MQTTListener      `yaml:"mqtt_listener,omitempty"`
	Processing        *Processing        `yaml:"processing,omitempty"`
	InfluxDBPublisher *InfluxDBPublisher `yaml:"influxdb_publisher"`
	TagNames          map[string]string  `yaml:"tag_names"`
}

func ReadConfig(configFile string) (Config, error) {
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("No config found! Tried to open \"%s\"\n", configFile)
		os.Exit(1)
	}

	f, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var conf Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)
	err = decoder.Decode(&conf)

	if err != nil {
		return Config{}, err
	}
	return conf, nil
}
