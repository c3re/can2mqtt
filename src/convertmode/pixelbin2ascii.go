package convertmode

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/brutella/can"
	"strconv"
	"strings"
)

type PixelBin2Ascii struct{}

func (_ PixelBin2Ascii) String() string {
	return "pixelbin2ascii"
}

func (_ PixelBin2Ascii) ToCan(input []byte) (can.Frame, error) {
	colorBytesAndNumber := strings.Fields(string(input))
	if len(colorBytesAndNumber) != 2 {
		return can.Frame{}, errors.New(fmt.Sprintf("input does not contain exactly two fields, one for the number and one for the color, got %d fields instead.", len(colorBytesAndNumber)))
	}
	colorBytes, _ := strings.CutPrefix(colorBytesAndNumber[1], "#")
	if len(colorBytes) != 6 {
		return can.Frame{}, errors.New(fmt.Sprintf("second field (color) does not contain exactly 6 nibbles each represented by one character, got %d instead", len(colorBytes)))
	}
	number, err := strconv.ParseUint(colorBytesAndNumber[0], 10, 8)
	if err != nil {
		return can.Frame{}, errors.New(fmt.Sprintf("Error while converting first field (pixel number): %s", err.Error()))
	}
	res, err := hex.DecodeString(colorBytes)
	if err != nil {
		return can.Frame{}, errors.New(fmt.Sprintf("Error while converting: %s", err.Error()))
	}
	var returner = can.Frame{Length: 4}
	returner.Data[0] = uint8(number)
	copy(res, returner.Data[1:4])
	return returner, nil
}

func (_ PixelBin2Ascii) ToMqtt(input can.Frame) ([]byte, error) {
	if input.Length != 4 {
		return []byte{}, errors.New(fmt.Sprintf("Input does not contain exactly 4 bytes, got %d instead", input.Length))
	}
	colorString := "#" + hex.EncodeToString(input.Data[0:3])
	numberString := strconv.FormatUint(uint64(input.Data[0]), 10)
	return []byte(strings.Join([]string{numberString, colorString}, " ")), nil
}