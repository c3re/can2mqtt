package convertfunctions

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/brutella/can"
	"strconv"
	"strings"
)

// Int2Ascii is a convertMode that can take multiple signed integers of one size.
// instances describe the amount of numbers that should be converted, bits is the size of each number
// instances * bits must fit into 64 bits.
type Int2Ascii struct {
	instances, bits uint
}

func (i2a Int2Ascii) String() string {
	instanceString := ""
	if i2a.instances > 1 {
		instanceString = fmt.Sprintf("%d", i2a.instances)
	}
	return fmt.Sprintf("%sint%d2ascii", instanceString, i2a.bits)
}

func NewInt2Ascii(instances, bits uint) (Int2Ascii, error) {
	if !(bits == 8 || bits == 16 || bits == 32 || bits == 64) {
		return Int2Ascii{}, errors.New(fmt.Sprintf("bitsize %d not supported, please choose one of 8, 16. 32 or 64\n", bits))
	}
	if bits*instances > 64 {
		return Int2Ascii{}, errors.New(fmt.Sprintf("%d instances of %d bit size would not fit into a 8 byte CAN-Frame. %d exceeds 64 bits.\n", instances, bits, instances*bits))
	}
	return Int2Ascii{instances, bits}, nil
}

// ToCan is the generic approach to convert instances instances of numbers with bits bits size.
// Allowed values for instances are 1-8.
// Allowed values for bits are 8, 16, 32 or 64
// instances*bits must not be larger than 64
// input has to contain the data that shall be converted. The input is split at whitespaces, the amount of fields has
// to match instances.
// If the amount of fields matches, each field is converted to an uint of size bits. The results are then added to the CAN-frame.
func (i2a Int2Ascii) ToCan(input []byte) (can.Frame, error) {
	splitInput := strings.Fields(string(input))
	if uint(len(splitInput)) != i2a.instances {
		return can.Frame{}, errors.New(fmt.Sprintf("input does not contain exactly %d numbers seperated by whitespace", i2a.instances))
	}
	var ret can.Frame
	ret.Length = uint8((i2a.instances * i2a.bits) >> 3)
	bytePerNumber := i2a.bits >> 3
	for i := uint(0); i < i2a.instances; i++ {
		res, err := strconv.ParseInt(splitInput[i], 10, int(i2a.bits))
		if err != nil {
			return can.Frame{}, errors.New(fmt.Sprintf("Error while converting string %d: %s, %s", i, splitInput[i], err))
		}
		switch i2a.bits {
		case 64:
			binary.LittleEndian.PutUint64(ret.Data[i*bytePerNumber:(i+1)*bytePerNumber], uint64(res))
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

// ToMqtt is the generic approach to convert instances instances of numbers with bits bits size.
// Allowed values for instances are 1-8.
// Allowed values for bits are 8, 16, 32 or 64
// instances*bits must not be larger than 64
// input has to Contain the Data that shall be converted. The Size of the CAN-Frame has to fit the expected size.
// If we have for example 1 amount of 32-Bits numbers the CAN-Frame size input.Length has to be 4 (bytes).
// If the size fits, the Data is split up in instances pieces and are then processed to a string representation
// via strconv.FormatUint.
// The successful return value is a byte-slice that represents the converted strings joined with a space between them.
func (i2a Int2Ascii) ToMqtt(input can.Frame) ([]byte, error) {
	if input.Length != uint8((i2a.bits*i2a.instances)>>3) {
		return []byte{}, errors.New(fmt.Sprintf("Input is of wrong length: %d, expected %d because of %d numbers of %d-bits.", input.Length, (i2a.instances*i2a.bits)>>3, i2a.instances, i2a.bits))
	}
	var returnStrings []string
	bytePerNumber := i2a.bits >> 3
	for i := uint(0); i < i2a.instances; i++ {
		switch i2a.bits {
		case 64:
			returnStrings = append(returnStrings, strconv.FormatInt(int64(binary.LittleEndian.Uint64(input.Data[i*bytePerNumber:(i+1)*bytePerNumber])), 10))
		case 32:
			returnStrings = append(returnStrings, strconv.FormatInt(int64(binary.LittleEndian.Uint32(input.Data[i*bytePerNumber:(i+1)*bytePerNumber])), 10))
		case 16:
			returnStrings = append(returnStrings, strconv.FormatInt(int64(binary.LittleEndian.Uint16(input.Data[i*bytePerNumber:(i+1)*bytePerNumber])), 10))
		case 8:
			returnStrings = append(returnStrings, strconv.FormatInt(int64(input.Data[i]), 10))
		}
	}
	return []byte(strings.Join(returnStrings, " ")), nil
}