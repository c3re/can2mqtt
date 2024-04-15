package convertfunctions

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

//######################################################################
//#				NONE				       #
//######################################################################

func bytes2ascii(length uint32, payload [8]byte) string {
	return string(payload[:length])
}

func ascii2bytes(payload string) ([8]byte, uint8) {
	var returner [8]byte
	var i uint8 = 0
	for ; int(i) < len(payload) && i < 8; i++ {
		returner[i] = payload[i]
	}
	return returner, i
}

// ######################################################################
// #			BOOL2ASCII				       #
// ######################################################################
// bool2ascii takes exactly two byte and returns a string with a
// boolean interpretation of the found data
func bool2ascii(payload []byte) string {
	if len(payload) != 2 {
		return "Err in CAN-Frame, data must be 2 bytes."
	}
	data := binary.LittleEndian.Uint16(payload)
	bits := strconv.FormatUint(uint64(data), 2)
	split := strings.Split(bits, "")
	// fill the '0' bits
	if len(split) < 16 {
		for i := len(split); i < 16; i++ {
			split = append([]string{"0"}, split...)
		}
	}
	// get the two 'bytes'
	lower := split[8:16]
	upper := split[0:8]
	// swap 'bytes', according to integer representation
	lower[0], lower[1], lower[2], lower[3], lower[4], lower[5], lower[6], lower[7] = lower[7], lower[6], lower[5], lower[4], lower[3], lower[2], lower[1], lower[0]
	upper[0], upper[1], upper[2], upper[3], upper[4], upper[5], upper[6], upper[7] = upper[7], upper[6], upper[5], upper[4], upper[3], upper[2], upper[1], upper[0]
	return strings.Join(lower, " ") + " " + strings.Join(upper, " ")
}

func ascii2bool(payload string) []byte {
	// split the 16 '0' or '1'
	split := strings.Split(payload, " ")
	// get the two 'bytes'
	lower := split[0:8]
	upper := split[8:16]
	// swap 'bytes', according to integer representation
	lower[0], lower[1], lower[2], lower[3], lower[4], lower[5], lower[6], lower[7] = lower[7], lower[6], lower[5], lower[4], lower[3], lower[2], lower[1], lower[0]
	upper[0], upper[1], upper[2], upper[3], upper[4], upper[5], upper[6], upper[7] = upper[7], upper[6], upper[5], upper[4], upper[3], upper[2], upper[1], upper[0]
	// convert to string again
	tmp := strings.Join(upper, "") + strings.Join(lower, "")
	number, err := strconv.ParseUint(tmp, 2, 16)
	a := make([]byte, 2)
	if err != nil {
		fmt.Printf("error converting %s\n", tmp)
	} else {
		binary.LittleEndian.PutUint16(a, uint16(number))
	}
	return a
}

// ######################################################################
// #			UINT82ASCII				       #
// ######################################################################
// uint82ascii takes exactly one byte and returns a string with a
// numeric decimal interpretation of the found data
func uint82ascii(payload byte) string {
	return strconv.FormatInt(int64(payload), 10)
}

func ascii2uint8(payload string) byte {
	return ascii2uint16(payload)[0]
}

// ######################################################################
// #			UINT162ASCII				       #
// ######################################################################
// uint162ascii takes 2 bytes and returns a string with a numeric
// decimal interpretation of the found data as ascii-string
func uint162ascii(payload []byte) string {
	if len(payload) != 2 {
		return "Err in CAN-Frame, data must be 2 bytes."
	}
	data := binary.LittleEndian.Uint16(payload)
	return strconv.FormatUint(uint64(data), 10)
}

func ascii2uint16(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint16(tmp)
	a := make([]byte, 2)
	binary.LittleEndian.PutUint16(a, number)
	return a
}

// ######################################################################
// #			INT162ASCII				       #
// ######################################################################
// int162ascii takes 2 bytes and returns a string with a numeric
// decimal interpretation of the found data as ascii-string
func int162ascii(payload []byte) string {
	if len(payload) != 2 {
		return "Err in CAN-Frame, data must be 2 bytes."
	}
	data := int16(binary.LittleEndian.Uint16(payload))
	return strconv.FormatInt(int64(data), 10)
}

func ascii2int16(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint16(tmp)
	a := make([]byte, 2)
	binary.LittleEndian.PutUint16(a, number)
	return a
}

// ########################################################################
// ######################################################################
// #			UINT322ASCII				       #
// ######################################################################
// uint322ascii takes 4 bytes and returns a string with a numeric
// decimal interpretation of the found data as ascii-string
func uint322ascii(payload []byte) string {
	if len(payload) != 4 {
		return "Err in CAN-Frame, data must be 4 bytes."
	}
	data := binary.LittleEndian.Uint32(payload)
	return strconv.FormatUint(uint64(data), 10)
}

func ascii2uint32(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint32(tmp)
	a := make([]byte, 4)
	binary.LittleEndian.PutUint32(a, number)
	return a
}

// ########################################################################
// ######################################################################
// #			UINT642ASCII				       #
// ######################################################################
// uint642ascii takes 8 bytes and returns a string with a numeric
// decimal interpretation of the found data as ascii-string
func uint642ascii(payload []byte) string {
	if len(payload) != 8 {
		return "Err in CAN-Frame, data must be 8 bytes."
	}
	data := binary.LittleEndian.Uint64(payload)
	return strconv.FormatUint(data, 10)
}

func ascii2uint64(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint64(tmp)
	a := make([]byte, 8)
	binary.LittleEndian.PutUint64(a, number)
	return a
}

// ########################################################################
// ######################################################################
// #             bytecolor2colorcode
// ######################################################################
// bytecolor2colorcode is a convertmode that converts between the binary
// 3 byte representation of a color and a string representation of a color
// as we know it (for example in html #00ff00 is green)
func bytecolor2colorcode(payload []byte) string {
	colorstring := hex.EncodeToString(payload)
	return "#" + colorstring
}

func colorcode2bytecolor(payload string) []byte {
	var a []byte
	var err error
	a, err = hex.DecodeString(strings.Replace(payload, "#", "", -1))
	if err != nil {
		return []byte{0, 0, 0}
	}
	return a
}

//########################################################################
