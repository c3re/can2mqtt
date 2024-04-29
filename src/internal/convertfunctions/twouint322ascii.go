package convertfunctions

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/brutella/can"
	"strconv"
	"strings"
)

func TwoUint322AsciiToCan(input []byte) (can.Frame, error) {
	splitInput := strings.Split(string(input), " ")
	if len(splitInput) != 2 {
		return can.Frame{}, errors.New("input does not contain exactly 2 numbers seperated by spaces")
	}
	var ret can.Frame
	ret.Length = 8
	for i := 0; i < 2; i++ {
		res, err := strconv.ParseUint(splitInput[i], 10, 32)
		if err != nil {
			return can.Frame{}, errors.New(fmt.Sprintf("Error while converting string %d: %s, %s", i, splitInput[i], err))
		}
		binary.LittleEndian.PutUint32(ret.Data[0+i*4:4+i*4], uint32(res))
	}
	return ret, nil
}

func TwoUint322AsciiToMqtt(input can.Frame) ([]byte, error) {
	if input.Length != 8 {
		return []byte{}, errors.New("input does not contain exactly 8 Bytes")
	}
	var returnStrings [2]string
	for i := 0; i < 2; i++ {
		returnStrings[i] = strconv.FormatUint(uint64(binary.LittleEndian.Uint32(input.Data[0+i*4:4+i*4])), 10)
	}
	return []byte(strings.Join(returnStrings[:], " ")), nil
}
