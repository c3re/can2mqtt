package main

import (
	"github.com/brutella/can"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
)

// handleCAN is the standard receive handler for CANFrames
// and does the following:
// 1. calling standard convert function: convert2MQTT
// 2. sending the message
func handleCAN(cf can.Frame) {
	slog.Debug("received CANFrame", "id", cf.ID, "len", cf.Length, "data", cf.Data)
	// Only do conversions when necessary
	if dirMode != MQTT2CAN_ONLY {
		mqttPayload, err := pairFromID[cf.ID].convertMode.ToMqtt(cf)
		if err != nil {
			slog.Warn("conversion to MQTT message unsuccessful", "convertmode", pairFromID[cf.ID].convertMode, "error", err)
			return
		}
		topic := pairFromID[cf.ID].mqttTopic
		mqttPublish(topic, mqttPayload)
		// this is the most common log-message, craft with care...
		slog.Info("CAN -> MQTT", "ID", cf.ID, "len", cf.Length, "data", cf.Data, "convertmode", pairFromID[cf.ID].convertMode, "topic", topic, "message", mqttPayload)
	}
}

// handleMQTT is the standard receive handler for MQTT
// messages and does the following:
// 1. calling the standard convert function: convert2CAN
// 2. sending the message
func handleMQTT(_ MQTT.Client, msg MQTT.Message) {
	slog.Debug("received message", "topic", msg.Topic(), "payload", msg.Payload())

	if dirMode != CAN2MQTT_ONLY {
		cf, err := pairFromTopic[msg.Topic()].convertMode.ToCan(msg.Payload())
		if err != nil {
			slog.Warn("conversion to CAN-Frame unsuccessful", "convertmode", pairFromTopic[msg.Topic()].convertMode, "error", err)
			return
		}
		cf.ID = pairFromTopic[msg.Topic()].canId
		canPublish(cf)
		slog.Info("CAN <- MQTT", "ID", cf.ID, "len", cf.Length, "data", cf.Data, "convertmode", pairFromTopic[msg.Topic()].convertMode, "topic", msg.Topic(), "message", msg.Payload())
	}
}
