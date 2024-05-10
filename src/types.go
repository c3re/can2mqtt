package main

import "github.com/brutella/can"

type convertToCan func(input []byte) (can.Frame, error)
type convertToMqtt func(input can.Frame) ([]byte, error)

// can2mqtt is a struct that represents the internal type of
// one line of the can2mqtt.csv file. It has
// the same three fields as the can2mqtt.csv file: CAN-ID,
// conversion method and MQTT-Topic.
type can2mqtt struct {
	canId      uint32
	convMethod string
	toCan      convertToCan
	toMqtt     convertToMqtt
	mqttTopic  string
}
