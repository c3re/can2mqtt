package internal

import (
	"fmt"
	"github.com/brutella/can"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// handleCAN is the standard receive handler for CANFrames
// and does the following:
// 1. calling standard convert function: convert2MQTT
// 2. sending the message
func handleCAN(cf can.Frame) {
	if dbg {
		fmt.Printf("receivehandler: received CANFrame: ID: %d, len: %d, payload %s\n", cf.ID, cf.Length, cf.Data)
	}
	// Only do conversions when necessary
	if dirMode != 2 {
		mqttPayload, err := pairFromID[cf.ID].toMqtt(cf)
		if err != nil {
			fmt.Printf("Error while converting CAN Frame with ID %d and payload %s: %s", cf.ID, cf.Data, err.Error())
			return
		}
		if dbg {
			fmt.Printf("receivehandler: converted String: %s\n", mqttPayload)
		}
		topic := getTopicFromId(cf.ID)
		mqttPublish(topic, mqttPayload)
		fmt.Printf("ID: %d len: %d data: %X -> topic: \"%s\" message: \"%s\"\n", cf.ID, cf.Length, cf.Data, topic, mqttPayload)
	}
}

// handleMQTT is the standard receive handler for MQTT
// messages and does the following:
// 1. calling the standard convert function: convert2CAN
// 2. sending the message
func handleMQTT(_ MQTT.Client, msg MQTT.Message) {
	if dbg {
		fmt.Printf("receivehandler: received message: topic: %s, msg: %s\n", msg.Topic(), msg.Payload())
	}

	if dirMode != 1 {
		//cf := convert2CAN(msg.Topic(), string(msg.Payload()))
		cf, err := pairFromTopic[msg.Topic()].toCan(msg.Payload())
		if err != nil {
			fmt.Printf("Error while converting MQTT-Message with Topic %s payload %s: %s", msg.Topic(), msg.Payload(), err.Error())
			return
		}
		if dbg {
			fmt.Printf("receivehandler: converted data: %s\n", cf.Data)
		}
		cf.ID = uint32(pairFromTopic[msg.Topic()].canId)
		canPublish(cf)
		fmt.Printf("ID: %d len: %d data: %X <- topic: \"%s\" message: \"%s\"\n", cf.ID, cf.Length, cf.Data, msg.Topic(), msg.Payload())
	}
}
