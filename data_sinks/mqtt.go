package data_sinks

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Scrin/RuuviBridge/common/limiter"
	"github.com/Scrin/RuuviBridge/config"
	"github.com/Scrin/RuuviBridge/parser"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

func MQTT(conf config.MQTTPublisher) chan<- parser.Measurement {
	address := conf.BrokerAddress
	if address == "" {
		address = "localhost"
	}
	port := conf.BrokerPort
	if port == 0 {
		port = 1883
	}
	server := conf.BrokerUrl
	if server == "" {
		server = fmt.Sprintf("tcp://%s:%d", address, port)
	}
	log.Info().
		Str("target", server).
		Str("topic_prefix", conf.TopicPrefix).
		Dur("minimum_interval", conf.MinimumInterval).
		Msg("Starting MQTT sink")

	clientID := conf.ClientID
	if clientID == "" {
		clientID = "RuuviBridgePublisher"
	}
	opts := mqtt.NewClientOptions()
	opts.SetCleanSession(false)
	opts.AddBroker(server)
	opts.SetClientID(clientID)
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	opts.SetKeepAlive(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	if conf.LWTTopic != "" {
		payload := conf.LWTOfflinePayload
		if payload == "" {
			payload = "{\"state\":\"offline\"}"
		}
		opts.SetWill(conf.LWTTopic, payload, 0, true)
	}
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Error().
			Str("target", server).
			Str("topic_prefix", conf.TopicPrefix).
			Dur("minimum_interval", conf.MinimumInterval).
			Err(token.Error()).
			Msg("Failed to connect to MQTT")
	}
	if conf.LWTTopic != "" {
		payload := conf.LWTOnlinePayload
		if payload == "" {
			payload = "{\"state\":\"online\"}"
		}
		client.Publish(conf.LWTTopic, 0, true, payload)
	}

	limiter := limiter.New(conf.MinimumInterval)
	measurements := make(chan parser.Measurement, 1024)
	go func() {
		for measurement := range measurements {
			if !limiter.Check(measurement) {
				log.Trace().Str("mac", measurement.Mac).Msg("Skipping MQTT publish due to interval limit")
				continue
			}
			data, err := json.Marshal(measurement)
			if err != nil {
				log.Error().Err(err).Msg("Failed to serialize measurement")
			} else {
				client.Publish(conf.TopicPrefix+"/"+measurement.Mac, 0, false, string(data))
				if conf.HomeassistantDiscoveryPrefix != "" {
					publishHomeAssistantDiscoveries(client, conf, measurement)
				}
				if conf.PublishRaw {
					safePublishF := func(label string, v *float64) {
						if v != nil {
							client.Publish(conf.TopicPrefix+"/"+measurement.Mac+"/"+label, 0, false, strconv.FormatFloat(*v, 'f', -1, 64))
						}
					}
					safePublishI := func(label string, v *int64) {
						if v != nil {
							client.Publish(conf.TopicPrefix+"/"+measurement.Mac+"/"+label, 0, false, strconv.FormatInt(*v, 10))
						}
					}
					safePublishB := func(label string, v *bool) {
						if v != nil {
							client.Publish(conf.TopicPrefix+"/"+measurement.Mac+"/"+label, 0, false, strconv.FormatBool(*v))
						}
					}
					safePublishF("temperature", measurement.Temperature)
					safePublishF("humidity", measurement.Humidity)
					safePublishF("pressure", measurement.Pressure)
					safePublishF("accelerationX", measurement.AccelerationX)
					safePublishF("accelerationY", measurement.AccelerationY)
					safePublishF("accelerationZ", measurement.AccelerationZ)
					safePublishF("batteryVoltage", measurement.BatteryVoltage)
					safePublishI("txPower", measurement.TxPower)
					safePublishI("rssi", measurement.Rssi)
					safePublishI("movementCounter", measurement.MovementCounter)
					safePublishI("measurementSequenceNumber", measurement.MeasurementSequenceNumber)
					safePublishF("accelerationTotal", measurement.AccelerationTotal)
					safePublishF("absoluteHumidity", measurement.AbsoluteHumidity)
					safePublishF("dewPoint", measurement.DewPoint)
					safePublishF("equilibriumVaporPressure", measurement.EquilibriumVaporPressure)
					safePublishF("airDensity", measurement.AirDensity)
					safePublishF("accelerationAngleFromX", measurement.AccelerationAngleFromX)
					safePublishF("accelerationAngleFromY", measurement.AccelerationAngleFromY)
					safePublishF("accelerationAngleFromZ", measurement.AccelerationAngleFromZ)
					// New E1 fields
					safePublishF("pm1p0", measurement.Pm1p0)
					safePublishF("pm2p5", measurement.Pm2p5)
					safePublishF("pm4p0", measurement.Pm4p0)
					safePublishF("pm10p0", measurement.Pm10p0)
					safePublishF("co2", measurement.CO2)
					safePublishF("voc", measurement.VOC)
					safePublishF("nox", measurement.NOX)
					safePublishF("illuminance", measurement.Illuminance)
					safePublishF("soundInstant", measurement.SoundInstant)
					safePublishF("soundAverage", measurement.SoundAverage)
					safePublishF("soundPeak", measurement.SoundPeak)
					// Diagnostics
					safePublishB("calibrationInProgress", measurement.CalibrationInProgress)
					safePublishB("buttonPressedOnBoot", measurement.ButtonPressedOnBoot)
					safePublishB("rtcOnBoot", measurement.RtcOnBoot)
				}
			}
		}
	}()
	return measurements
}
