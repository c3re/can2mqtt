package convertmode

import "github.com/brutella/can"

type None struct{}

func (_ None) String() string {
	return "none"
}

func (_ None) ToCan(input []byte) (can.Frame, error) {
	var returner [8]byte
	var i uint8 = 0
	for ; int(i) < len(input) && i < 8; i++ {
		returner[i] = input[i]
	}
	return can.Frame{Length: i, Data: returner}, nil
}

func (_ None) ToMqtt(input can.Frame) ([]byte, error) {
	return input.Data[:input.Length], nil
}

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