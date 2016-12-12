package main

import (
	"fmt"
	CAN "github.com/brendoncarroll/go-can"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// handleCAN is the standard receivehandler for CANFrames
// and does the following:
// 1. calling standard convertfunction: convert2MQTT
// 2. sending the message
func handleCAN(cf CAN.CANFrame) {
	if dbg {
		fmt.Printf("receivehandler: received CANFrame: ID: %d, len: %d, payload %s\n", cf.ID, cf.Len, cf.Data)
	}
	mqttPayload := convert2MQTT(int(cf.ID), int(cf.Len), cf.Data)
	if dbg {
		fmt.Printf("receivehandler: converted String: %s\n", mqttPayload)
	}
	topic := getTopic(int(cf.ID))
	mqttPublish(topic, mqttPayload)
	fmt.Printf("ID: %d len: %d data: %X -> topic: \"%s\" message: \"%s\"\n", cf.ID, cf.Len, cf.Data, topic, mqttPayload)
}

// handleMQTT is the standard receivehandler for MQTT
// messages and does the following:
// 1. calling the standard convertfunction: convert2CAN
// 2. sending the message
func handleMQTT(cl MQTT.Client, msg MQTT.Message) {
	if dbg {
		fmt.Printf("receivehandler: received message: topic: %s, msg: %s\n", msg.Topic(), msg.Payload())
	}
	cf := convert2CAN(msg.Topic(), string(msg.Payload()))
	canPublish(cf)
	fmt.Printf("ID: %d len: %d data: %X <- topic: \"%s\" message: \"%s\"\n", cf.ID, cf.Len, cf.Data, msg.Topic(), msg.Payload())
}
