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
