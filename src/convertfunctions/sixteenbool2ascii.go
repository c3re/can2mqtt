package convertfunctions

import (
	"errors"
	"fmt"
	"github.com/brutella/can"
	"strconv"
	"strings"
)

func SixteenBool2AsciiToCan(input []byte) (can.Frame, error) {
	splitInput := strings.Split(string(input), " ") // TODO use strings.Fields here
	if len(splitInput) != 16 {
		return can.Frame{}, errors.New("input does not contain exactly 16 numbers seperated by spaces")
	}
	var returnData [8]uint8
	for i := 0; i < len(splitInput); i++ {
		res, err := strconv.ParseBool(splitInput[i])
		if err != nil {
			return can.Frame{}, errors.New(fmt.Sprintf("input does not specify a boolean at index %d: %s:%s", i, splitInput[i], err))
		}
		if res {
			returnData[i>>3] |= 0x1 << (i % 8)
		} else {

			returnData[i>>3] |= 0x0 << (i % 8)
		}
	}
	return can.Frame{Length: 2, Data: returnData}, nil
}
func SixteenBool2AsciiToMqtt(input can.Frame) ([]byte, error) {
	var returnStrings [16]string
	for i := 0; i < 16; i++ {
		if (input.Data[i>>3]>>(i%8))&0x1 == 1 {
			returnStrings[i] = "1"
		} else {
			returnStrings[i] = "0"
		}
	}
	return []byte(strings.Join(returnStrings[:], " ")), nil
}
