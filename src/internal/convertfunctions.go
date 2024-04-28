package internal

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/brutella/can"
	"strconv"
	"strings"
)

// convert2CAN does the following:
// 1. receive topic and payload
// 2. use topic to examine corresponding convertmode and CAN-ID
// 3. execute conversion
// 4. build CANFrame
// 5. returning the CANFrame
func convert2CAN(topic, payload string) can.Frame {
	//convertMethod := getConvModeFromTopic(topic)
	//var Id = uint32(getIdFromTopic(topic))
	//var data [8]byte
	//var length uint8
	frame, err := pairFromTopic[topic].toCan(payload)
	if err != nil {
		fmt.Printf("Error while converting %s\n", err.Error())
		return can.Frame{}
	}
	return frame
	/*
		if convertMethod == "none" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode none (reverse of %s)\n", convertMethod)
			}
			convertedFrame, _ := convert.NoneToCan(payload)
			// TODO check error
			data = convertedFrame.Data
			length = convertedFrame.Length
		} else if convertMethod == "16bool2ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii2bool (reverse of %s)\n", convertMethod)
			}
				convertedFrame, _ := convert.16Bool2AsciiToCan(payload)
				// TODO check error
				data = convertedFrame.Data
				length = convertedFrame.Length

			tmp := ascii2bool(payload)
			data[0] = tmp[0]
			data[1] = tmp[1]
			length = 2
		} else if convertMethod == "uint82ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii2uint8 (reverse of %s)\n", convertMethod)
			}
			data[0] = ascii2uint8(payload)
			length = 1
		} else if convertMethod == "uint162ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii2uint16(reverse of %s)\n", convertMethod)
			}
			tmp := ascii2uint16(payload)
			data[0] = tmp[0]
			data[1] = tmp[1]
			length = 2
		} else if convertMethod == "uint322ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii2uint32(reverse of %s)\n", convertMethod)
			}
			tmp := ascii2uint32(payload)
			data[0] = tmp[0]
			data[1] = tmp[1]
			data[2] = tmp[2]
			data[3] = tmp[3]
			length = 4
		} else if convertMethod == "uint642ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii2uint64(reverse of %s)\n", convertMethod)
			}
			tmp := ascii2uint64(payload)
			data[0] = tmp[0]
			data[1] = tmp[1]
			data[2] = tmp[2]
			data[3] = tmp[3]
			data[4] = tmp[4]
			data[5] = tmp[5]
			data[6] = tmp[6]
			data[7] = tmp[7]
			length = 8
		} else if convertMethod == "2uint322ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii22uint32(reverse of %s)\n", convertMethod)
			}
			nums := strings.Split(payload, " ")
			tmp := ascii2uint32(nums[0])
			data[0] = tmp[0]
			data[1] = tmp[1]
			data[2] = tmp[2]
			data[3] = tmp[3]
			tmp = ascii2uint32(nums[1])
			data[4] = tmp[0]
			data[5] = tmp[1]
			data[6] = tmp[2]
			data[7] = tmp[3]
			length = 8
		} else if convertMethod == "4uint162ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii24uint16(reverse of %s)\n", convertMethod)
			}
			nums := strings.Split(payload, " ")
			tmp := ascii2uint16(nums[0])
			data[0] = tmp[0]
			data[1] = tmp[1]
			tmp = ascii2uint16(nums[1])
			data[2] = tmp[0]
			data[3] = tmp[1]
			tmp = ascii2uint16(nums[2])
			data[4] = tmp[0]
			data[5] = tmp[1]
			tmp = ascii2uint16(nums[3])
			data[6] = tmp[0]
			data[7] = tmp[1]
			length = 8
		} else if convertMethod == "4int162ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii24int16(reverse of %s)\n", convertMethod)
			}
			nums := strings.Split(payload, " ")
			tmp := ascii2int16(nums[0])
			data[0] = tmp[0]
			data[1] = tmp[1]
			tmp = ascii2int16(nums[1])
			data[2] = tmp[0]
			data[3] = tmp[1]
			tmp = ascii2int16(nums[2])
			data[4] = tmp[0]
			data[5] = tmp[1]
			tmp = ascii2int16(nums[3])
			data[6] = tmp[0]
			data[7] = tmp[1]
			length = 8
		} else if convertMethod == "4uint82ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii24uint8(reverse of %s)\n", convertMethod)
			}
			nums := strings.Split(payload, " ")
			tmp := ascii2uint8(nums[0])
			data[0] = tmp
			data[1] = 0
			tmp = ascii2uint8(nums[1])
			data[2] = tmp
			data[3] = 0
			tmp = ascii2uint8(nums[2])
			data[4] = tmp
			data[5] = 0
			tmp = ascii2uint8(nums[3])
			data[6] = tmp
			data[7] = 0
			length = 8
		} else if convertMethod == "8uint82ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii28uint8(reverse of %s)\n", convertMethod)
			}
			nums := strings.Split(payload, " ")
			if len(nums) != 8 {
				fmt.Printf("Error, wrong number of bytes provided for convertmode 8uint82ascii, expected 8 got %d\n", len(nums))
			}
			data[0] = ascii2uint8(nums[0])
			data[1] = ascii2uint8(nums[1])
			data[2] = ascii2uint8(nums[2])
			data[3] = ascii2uint8(nums[3])
			data[4] = ascii2uint8(nums[4])
			data[5] = ascii2uint8(nums[5])
			data[6] = ascii2uint8(nums[6])
			data[7] = ascii2uint8(nums[7])
			length = 8
		} else if convertMethod == "bytecolor2colorcode" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode colorcode2bytecolor(reverse of %s)\n", convertMethod)
			}
			tmp := colorcode2bytecolor(payload)
			data[0] = tmp[0]
			data[1] = tmp[1]
			data[2] = tmp[2]
			length = 3
		} else if convertMethod == "pixelbin2ascii" {
			if dbg {
				fmt.Printf("convertfunctions: using convertmode ascii2pixelbin(reverse of %s)\n", convertMethod)
			}
			numAndColor := strings.Split(payload, " ")
			binNum := ascii2uint8(numAndColor[0])
			tmp := colorcode2bytecolor(numAndColor[1])
			data[0] = binNum
			data[1] = tmp[0]
			data[2] = tmp[1]
			data[3] = tmp[2]
			length = 4
		} else {
			if dbg {
				fmt.Printf("convertfunctions: convertmode %s not found. using fallback none\n", convertMethod)
			}
			data, length = ascii2bytes(payload)
		}
		myFrame := can.Frame{ID: Id, Length: length, Data: data}
		return myFrame
	*/
}

// convert2MQTT does the following
// 1. receive ID and payload
// 2. lookup the correct convertmode
// 3. executing conversion
// 4. building a string
// 5. return
func convert2MQTT(id int, length int, payload [8]byte) string {
	convertMethod := getConvModeFromId(id)
	if convertMethod == "none" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode none\n")
		}
		return bytes2ascii(uint32(length), payload)
	} else if convertMethod == "16bool2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode 16bool2ascii\n")
		}
		return bool2ascii(payload[0:2])
	} else if convertMethod == "uint82ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode uint82ascii\n")
		}
		return uint82ascii(payload[0])
	} else if convertMethod == "uint162ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode uint162ascii\n")
		}
		return uint162ascii(payload[0:2])
	} else if convertMethod == "uint322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode uint322ascii\n")
		}
		return uint322ascii(payload[0:4])
	} else if convertMethod == "uint642ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode uint642ascii\n")
		}
		return uint642ascii(payload[0:8])
	} else if convertMethod == "2uint322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode 2uint322ascii\n")
		}
		return uint322ascii(payload[0:4]) + " " + uint322ascii(payload[4:8])
	} else if convertMethod == "4uint162ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode 4uint162ascii\n")
		}
		return uint162ascii(payload[0:2]) + " " + uint162ascii(payload[2:4]) + " " + uint162ascii(payload[4:6]) + " " + uint162ascii(payload[6:8])
	} else if convertMethod == "4int162ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode 4int162ascii\n")
		}
		return int162ascii(payload[0:2]) + " " + int162ascii(payload[2:4]) + " " + int162ascii(payload[4:6]) + " " + int162ascii(payload[6:8])
	} else if convertMethod == "4uint82ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode 4uint82ascii\n")
		}
		return uint82ascii(payload[0]) + " " + uint82ascii(payload[2]) + " " + uint82ascii(payload[4]) + " " + uint82ascii(payload[6])
	} else if convertMethod == "8uint82ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode 8uint82ascii\n")
		}
		return uint82ascii(payload[0]) + " " + uint82ascii(payload[1]) + " " + uint82ascii(payload[2]) + " " + uint82ascii(payload[3]) + " " + uint82ascii(payload[4]) + " " + uint82ascii(payload[5]) + " " + uint82ascii(payload[6]) + " " + uint82ascii(payload[7])
	} else if convertMethod == "pixelbin2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode pixelbin2ascii\n")
		}
		return uint82ascii(payload[0]) + " " + bytecolor2colorcode(payload[1:4])
	} else if convertMethod == "bytecolor2colorcode" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode bytecolor2colorcode\n")
		}
		return bytecolor2colorcode(payload[0:2])

	} else {
		if dbg {
			fmt.Printf("convertfunctions: convertmode %s not found. using fallback none\n", convertMethod)
		}
		return bytes2ascii(uint32(length), payload)
	}
}

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
