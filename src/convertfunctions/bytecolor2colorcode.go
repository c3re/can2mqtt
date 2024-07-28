package convertfunctions

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/brutella/can"
	"strings"
)

type ByteColor2ColorCode struct{}

func (_ ByteColor2ColorCode) String() string {
	return "bytecolor2colorcode"
}

func (_ ByteColor2ColorCode) ToCan(input []byte) (can.Frame, error) {
	colorBytes, _ := strings.CutPrefix(string(input), "#")
	if len(colorBytes) != 6 {
		return can.Frame{}, errors.New(fmt.Sprintf("input does not contain exactly 6 nibbles each represented by one character, got %d instead", len(colorBytes)))
	}
	res, err := hex.DecodeString(colorBytes)
	if err != nil {
		return can.Frame{}, errors.New(fmt.Sprintf("Error while converting: %s", err.Error()))
	}
	var returner = can.Frame{Length: 3}
	copy(res, returner.Data[0:3])
	return returner, nil
}

func (_ ByteColor2ColorCode) ToMqtt(input can.Frame) ([]byte, error) {
	if input.Length != 3 {
		return []byte{}, errors.New(fmt.Sprintf("Input does not contain exactly 3 bytes, got %d instead", input.Length))
	}
	colorstring := hex.EncodeToString(input.Data[0:3])
	return []byte("#" + colorstring), nil
}