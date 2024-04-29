package convertfunctions

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/brutella/can"
	"strconv"
)

func Uint82AsciiToCan(input []byte) (can.Frame, error) {
	return UintN2AsciiToCan(8, input)
}

func Uint82AsciiToMqtt(input can.Frame) ([]byte, error) {
	return UintN2AsciiToMqtt(8, input)
}

func Uint162AsciiToCan(input []byte) (can.Frame, error) {
	return UintN2AsciiToCan(16, input)
}

func Uint162AsciiToMqtt(input can.Frame) ([]byte, error) {
	return UintN2AsciiToMqtt(16, input)
}

func Uint322AsciiToCan(input []byte) (can.Frame, error) {
	return UintN2AsciiToCan(32, input)
}

func Uint322AsciiToMqtt(input can.Frame) ([]byte, error) {
	return UintN2AsciiToMqtt(32, input)
}

func Uint642AsciiToCan(input []byte) (can.Frame, error) {
	return UintN2AsciiToCan(64, input)
}

func Uint642AsciiToMqtt(input can.Frame) ([]byte, error) {
	return UintN2AsciiToMqtt(64, input)
}

func UintN2AsciiToMqtt(n uint8, input can.Frame) ([]byte, error) {
	if input.Length != (n >> 3) {
		return []byte{}, errors.New(fmt.Sprintf("Error converting to MQTT, expected exactly %d byte got %d", n>>3, input.Length))
	}
	var tmpUint64 uint64
	for i := uint8(0); i < n>>3; i++ {
		tmpUint64 |= uint64(input.Data[i]) << (i << 3)
	}
	return []byte(strconv.FormatUint(tmpUint64, 10)), nil
}
func UintN2AsciiToCan(n uint8, input []byte) (can.Frame, error) {
	res, err := strconv.ParseUint(string(input), 10, int(n))
	if err != nil {
		return can.Frame{}, errors.New(fmt.Sprintf("Error converting to CAN: %s", err))
	}
	var ret can.Frame
	binary.LittleEndian.PutUint64(ret.Data[:], res)
	ret.Length = n >> 3
	return ret, nil
}
