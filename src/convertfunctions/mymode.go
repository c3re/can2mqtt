package convertfunctions

import (
	"errors"
	"github.com/brutella/can"
)

const mockErr string = "I am just mockup-code and not supposed to be actually used, implement something useful here"

func MyModeToCan(input []byte) (can.Frame, error) {
	/*
		This is your area to create your convertMode (Receive MQTT, convert to CAN).
		You can find the payload of the received MQTT-Message
		in the []byte input. You can craft your returning can-Frame here. It does not make sense to set the ID,
		it will be overwritten. You can also return an error, the Frame is not sent in that case.

		As an example you could use the following code to implement the "none" convert-Mode.

		var returner [8]byte
		var i uint8 = 0
		for ; int(i) < len(input) && i < 8; i++ {
			returner[i] = input[i]
		}
		return can.Frame{Length: i, Data: returner}, nil
	*/
	return can.Frame{}, errors.New(mockErr)
}

func MyModeToMqtt(input can.Frame) ([]byte, error) {
	/*
		This is your area to create your convertMode (Receive CAN, convert to MQTT).
		You can find the received CAN-Frame in the can.Frame input. You can craft your returning MQTT-Payload here.
		You can also return an error, the MQTT-Message is not sent in that case.

		As an example you could use the following code to implement the "none" convert-Mode.

		return input.Data[:input.Length], nil
	*/
	return []byte{}, errors.New(mockErr)
}
