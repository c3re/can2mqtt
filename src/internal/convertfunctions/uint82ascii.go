package convertfunctions

import (
	"errors"
	"fmt"
	"github.com/brutella/can"
	"strconv"
)

func Uint82AsciiToCan(input []byte) (can.Frame, error) {
	if len(input) != 1 {
		return can.Frame{}, errors.New(fmt.Sprintf("Error converting to CAN, expected exactly 1 byte got %d", len(input)))
	}
	res, err := strconv.ParseUint(string(input), 10, 8)
	if err != nil {
		return can.Frame{}, errors.New(fmt.Sprintf("Error converting to CAN: %s", err))
	}
	return can.Frame{Length: 1, Data: [8]uint8{uint8(res)}}, nil
}

func Uint82AsciiToMqtt(input can.Frame) ([]byte, error) {
	if input.Length != 1 {
		return []byte{}, errors.New(fmt.Sprintf("Error converting to MQTT, expected exactly 1 byte got %d", input.Length))
	}
	return []byte(strconv.FormatUint(uint64(input.Data[0]), 10)), nil
}
