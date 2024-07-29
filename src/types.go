package main

import (
	"fmt"
	"github.com/brutella/can"
)

// ConvertMode is the interface that defines the two methods necessary
// to handle MQTT-Messages (ToMqtt) as well as CAN-Frames(ToCan). It also includes fmt.Stringer
// to make types that implement it print their human-readable convertmode, as it
// appears in the can2mqtt file.
type ConvertMode interface {
	ToCan(input []byte) (can.Frame, error)
	ToMqtt(input can.Frame) ([]byte, error)
	fmt.Stringer
}

// can2mqtt is a struct that represents the internal type of
// one line of the can2mqtt.csv file. It has
// the same three fields as the can2mqtt.csv file: CAN-ID,
// conversion method and MQTT-Topic.
type can2mqtt struct {
	canId       uint32
	convertMode ConvertMode
	mqttTopic   string
}

// Valid values for "dirMode"
const (
	BIDIRECTIONAL = iota
	CAN2MQTT_ONLY
	MQTT2CAN_ONLY
)