package convertfunctions

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/brutella/can"
	"strconv"
	"strings"
)

func Uint82AsciiToCan(input []byte) (can.Frame, error) {
	return NUintM2AsciiToCan(1, 8, input)
}

func Uint82AsciiToMqtt(input can.Frame) ([]byte, error) {
	return NUintM2AsciiToMqtt(1, 8, input)
}

func Uint162AsciiToCan(input []byte) (can.Frame, error) {
	return NUintM2AsciiToCan(1, 16, input)
}

func Uint162AsciiToMqtt(input can.Frame) ([]byte, error) {
	return NUintM2AsciiToMqtt(1, 16, input)
}

func Uint322AsciiToCan(input []byte) (can.Frame, error) {
	return NUintM2AsciiToCan(1, 32, input)
}

func Uint322AsciiToMqtt(input can.Frame) ([]byte, error) {
	return NUintM2AsciiToMqtt(1, 32, input)
}

func Uint642AsciiToCan(input []byte) (can.Frame, error) {
	return NUintM2AsciiToCan(1, 64, input)
}

func Uint642AsciiToMqtt(input can.Frame) ([]byte, error) {
	return NUintM2AsciiToMqtt(1, 64, input)
}

func TwoUint322AsciiToCan(input []byte) (can.Frame, error) {
	return NUintM2AsciiToCan(2, 32, input)
}

func TwoUint322AsciiToMqtt(input can.Frame) ([]byte, error) {
	return NUintM2AsciiToMqtt(2, 32, input)
}

func EightUint82AsciiToCan(input []byte) (can.Frame, error) {
	return NUintM2AsciiToCan(8, 8, input)
}

func EightUint82AsciiToMqtt(input can.Frame) ([]byte, error) {
	return NUintM2AsciiToMqtt(8, 8, input)
}

func FourUint82AsciiToCan(input []byte) (can.Frame, error) {
	return NUintM2AsciiToCan(4, 8, input)
}

func FourUint82AsciiToMqtt(input can.Frame) ([]byte, error) {
	return NUintM2AsciiToMqtt(4, 8, input)
}

func FourUint162AsciiToCan(input []byte) (can.Frame, error) {
	return NUintM2AsciiToCan(4, 16, input)
}

func FourUint162AsciiToMqtt(input can.Frame) ([]byte, error) {
	return NUintM2AsciiToMqtt(4, 16, input)
}

// NUintM2AsciiToCan is the generic approach to convert numberAmount occurrences of numbers with numberWidth bits size.
// Allowed values for numberAmount are 1-8.
// Allowed values for numberWidth are 8, 16, 32 or 64
// numberAmount*numberWidth shall not be larger than 64
// input has to contain the data that shall be converted. The input is split at whitespaces, the amount of fields has
// to match numberAmount.
// If the amount of fields matches, each field is converted to an uint of size numberWidth. The results are then added to the CAN-frame.
func NUintM2AsciiToCan(numberAmount, numberWidth uint, input []byte) (can.Frame, error) {
	if !(numberWidth == 8 || numberWidth == 16 || numberWidth == 32 || numberWidth == 64) {

		return can.Frame{}, errors.New(fmt.Sprintf("numberWitdh %d uknown please choose one of 8, 16. 32 or 64\n", numberWidth))

	}
	if numberWidth*numberAmount > 64 {
		return can.Frame{}, errors.New(fmt.Sprintf("%d number of %d bit width would not fit into a 8 byte CAN-Frame %d exceeds 64 bits.\n", numberAmount, numberWidth, numberAmount*numberWidth))
	}
	splitInput := strings.Fields(string(input))
	if uint(len(splitInput)) != numberAmount {
		return can.Frame{}, errors.New(fmt.Sprintf("input does not contain exactly %d numbers seperated by whitespace", numberAmount))
	}
	var ret can.Frame
	ret.Length = uint8((numberAmount * numberWidth) >> 3)
	bytePerNumber := numberWidth >> 3
	for i := uint(0); i < numberAmount; i++ {
		res, err := strconv.ParseUint(splitInput[i], 10, int(numberWidth))
		if err != nil {
			return can.Frame{}, errors.New(fmt.Sprintf("Error while converting string %d: %s, %s", i, splitInput[i], err))
		}
		switch numberWidth {
		case 64:
			binary.LittleEndian.PutUint64(ret.Data[i*bytePerNumber:(i+1)*bytePerNumber], res)
		case 32:
			binary.LittleEndian.PutUint32(ret.Data[i*bytePerNumber:(i+1)*bytePerNumber], uint32(res))
		case 16:
			binary.LittleEndian.PutUint16(ret.Data[i*bytePerNumber:(i+1)*bytePerNumber], uint16(res))
		case 8:
			ret.Data[i] = uint8(res)
		}
	}
	return ret, nil
}

// NUintM2AsciiToMqtt is the generic approach to convert numberAmount occurrences of numbers with numberWidth bits size.
// Allowed values for numberAmount are 1-8.
// Allowed values for numberWidth are 8, 16, 32 or 64
// numberAmount*numberWidth shall not be larger than 64
// input has to Contain the Data that shall be converted. The Size of the CAN-Frame has to fit the expected size.
// If we have for example 1 amount of 32-Bits numbers the CAN-Frame size input.Length has to be 4 (bytes).
// If the size fits, the Data is split up in numberAmount pieces and are then processed to a string representation
// via strconv.FormatUint.
// The successful return value is a byte-slice that represents the converted strings joined with a space between them.
func NUintM2AsciiToMqtt(numberAmount, numberWidth uint, input can.Frame) ([]byte, error) {
	if !(numberWidth == 8 || numberWidth == 16 || numberWidth == 32 || numberWidth == 64) {
		return []byte{}, errors.New(fmt.Sprintf("numberWitdh %d uknown please choose one of 8, 16. 32 or 64\n", numberWidth))
	}
	if numberWidth*numberAmount > 64 {
		return []byte{}, errors.New(fmt.Sprintf("%d number of %d bit width would not fit into a 8 byte CAN-Frame %d exceeds 64 bits.\n", numberAmount, numberWidth, numberAmount*numberWidth))
	}
	if input.Length != uint8((numberWidth*numberAmount)>>3) {
		return []byte{}, errors.New(fmt.Sprintf("Input is of wrong length: %d, expected %d because of %d numbers of %d-bits.", input.Length, (numberAmount*numberWidth)>>3, numberAmount, numberWidth))
	}
	var returnStrings []string
	bytePerNumber := numberWidth >> 3
	for i := uint(0); i < numberAmount; i++ {
		switch numberWidth {
		case 64:
			returnStrings = append(returnStrings, strconv.FormatUint(binary.LittleEndian.Uint64(input.Data[i*bytePerNumber:(i+1)*bytePerNumber]), 10))
		case 32:
			returnStrings = append(returnStrings, strconv.FormatUint(uint64(binary.LittleEndian.Uint32(input.Data[i*bytePerNumber:(i+1)*bytePerNumber])), 10))
		case 16:
			returnStrings = append(returnStrings, strconv.FormatUint(uint64(binary.LittleEndian.Uint16(input.Data[i*bytePerNumber:(i+1)*bytePerNumber])), 10))
		case 8:
			returnStrings = append(returnStrings, strconv.FormatUint(uint64(input.Data[i]), 10))
		}
	}
	return []byte(strings.Join(returnStrings, " ")), nil
}
