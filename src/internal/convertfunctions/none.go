package convertfunctions

import "github.com/brutella/can"

func NoneToCan(input []byte) (can.Frame, error) {
	var returner [8]byte
	var i uint8 = 0
	for ; int(i) < len(input) && i < 8; i++ {
		returner[i] = input[i]
	}
	return can.Frame{Length: i, Data: returner}, nil
}

func NoneToMqtt(input can.Frame) ([]byte, error) {
	return input.Data[:input.Length], nil
}
