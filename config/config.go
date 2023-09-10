package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type GatewayPolling struct {
	Enabled     *bool         `yaml:"enabled,omitempty"`
	GatewayUrl  string        `yaml:"gateway_url"`
	BearerToken string        `yaml:"bearer_token"`
	Interval    time.Duration `yaml:"interval"`
}

type MQTTListener struct {
	Enabled           *bool  `yaml:"enabled,omitempty"`
	BrokerUrl         string `yaml:"broker_url"`
	BrokerAddress     string `yaml:"broker_address"`
	BrokerPort        int    `yaml:"broker_port"`
	ClientID          string `yaml:"client_id"`
	Username          string `yaml:"username"`
	Password          string `yaml:"password"`
	TopicPrefix       string `yaml:"topic_prefix"`
	LWTTopic          string `yaml:"lwt_topic"`
	LWTOnlinePayload  string `yaml:"lwt_online_payload"`
	LWTOfflinePayload string `yaml:"lwt_offline_payload"`
}

type HTTPListener struct {
	Enabled *bool `yaml:"enabled,omitempty"`
	Port    int   `yaml:"port"`
}

type Processing struct {
	ExtendedValues *bool    `yaml:"extended_values,omitempty"`
	FilterMode     string   `yaml:"filter_mode"`
	FilterList     []string `yaml:"filter_list"`
}

type InfluxDBPublisher struct {
	Enabled         *bool             `yaml:"enabled,omitempty"`
	MinimumInterval time.Duration     `yaml:"minimum_interval,omitempty"`
	Url             string            `yaml:"url"`
	AuthToken       string            `yaml:"auth_token"`
	Org             string            `yaml:"org"`
	Bucket          string            `yaml:"bucket"`
	Measurement     string            `yaml:"measurement"`
	AdditionalTags  map[string]string `yaml:"additional_tags,omitempty"`
}

type Prometheus struct {
	Enabled *bool `yaml:"enabled,omitempty"`
	Port    int   `yaml:"port"`
}

type MQTTPublisher struct {
	Enabled                      *bool         `yaml:"enabled,omitempty"`
	MinimumInterval              time.Duration `yaml:"minimum_interval,omitempty"`
	BrokerUrl                    string        `yaml:"broker_url"`
	BrokerAddress                string        `yaml:"broker_address"`
	BrokerPort                   int           `yaml:"broker_port"`
	ClientID                     string        `yaml:"client_id"`
	Username                     string        `yaml:"username"`
	Password                     string        `yaml:"password"`
	TopicPrefix                  string        `yaml:"topic_prefix"`
	PublishRaw                   bool          `yaml:"publish_raw"`
	HomeassistantDiscoveryPrefix string        `yaml:"homeassistant_discovery_prefix,omitempty"`
	LWTTopic                     string        `yaml:"lwt_topic"`
	LWTOnlinePayload             string        `yaml:"lwt_online_payload"`
	LWTOfflinePayload            string        `yaml:"lwt_offline_payload"`
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
	HTTPListener      *HTTPListener      `yaml:"http_listener,omitempty"`
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
