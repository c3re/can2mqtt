package can2mqtt

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	CAN "github.com/brendoncarroll/go-can"
	"strconv"
	"strings"
)

// convert2CAN does the following:
// 1. receive topic and payload
// 2. use topic to examine corresponding cconvertmode and CAN-ID
// 3. execute conversion
// 4. build CANFrame
// 5. returning the CANFrame
func convert2CAN(topic, payload string) CAN.CANFrame {
	convertMethod := getConvTopic(topic)
	var Id uint32 = uint32(getId(topic))
	var data [8]byte
	var len uint32
	if convertMethod == "none" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode none (reverse of %s)\n", convertMethod)
		}
		data, len = ascii2bytes(payload)
	} else if convertMethod == "uint82ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2uint8 (reverse of %s)\n", convertMethod)
		}
		data[0] = ascii2uint8(payload)
		len = 1
	} else if convertMethod == "uint162ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2uint16(reverse of %s)\n", convertMethod)
		}
		tmp := ascii2uint16(payload)
		data[0] = tmp[0]
		data[1] = tmp[1]
		len = 2
	} else if convertMethod == "uint322ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2uint32(reverse of %s)\n", convertMethod)
		}
		tmp := ascii2uint32(payload)
		data[0] = tmp[0]
		data[1] = tmp[1]
		data[2] = tmp[2]
		data[3] = tmp[3]
		len = 4
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
		len = 8
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
		len = 8
	} else if convertMethod == "bytecolor2colorcode" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode colorcode2bytecolor(reverse of %s)\n", convertMethod)
		}
		tmp := colorcode2bytecolor(payload)
		data[0] = tmp[0]
		data[1] = tmp[1]
		data[2] = tmp[2]
		len = 3
	} else if convertMethod == "pixelbin2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2pixelbin(reverse of %s)\n", convertMethod)
		}
		num_and_color := strings.Split(payload, " ")
		bin_num := ascii2uint8(num_and_color[0])
		tmp := colorcode2bytecolor(num_and_color[1])
		data[0] = byte(bin_num)
		data[1] = tmp[0]
		data[2] = tmp[1]
		data[3] = tmp[2]
		len = 4
	} else {
		if dbg {
			fmt.Printf("convertfunctions: convertmode %s not found. using fallback none\n", convertMethod)
		}
		data, len = ascii2bytes(payload)
	}
	mycf := CAN.CANFrame{ID: Id, Len: len, Data: data}
	return mycf
}

// convert2MQTT does the following
// 1. receive ID and payload
// 2. lookup the correct convertmode
// 3. executing conversion
// 4. building a string
// 5. return
func convert2MQTT(id uint32, length int, payload [8]byte) string {
	convertMethod := getConvId(id)
	if convertMethod == "none" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode none\n")
		}
		return bytes2ascii(uint32(length), payload)
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

func ascii2bytes(payload string) ([8]byte, uint32) {
	var returner [8]byte
	var i uint32 = 0
	for ; int(i) < len(payload) && i < 8; i++ {
		returner[i] = payload[i]
	}
	return returner, i
}

//######################################################################
//#			UINT82ASCII				       #
//######################################################################
// uint82ascii takes exactly one byte and returns a string with a
// numeric decimal interpretation of the found data
func uint82ascii(payload byte) string {
	return strconv.FormatInt(int64(payload), 10)
}

func ascii2uint8(payload string) byte {
	return ascii2uint16(payload)[0]
}

//######################################################################
//#			UINT162ASCII				       #
//######################################################################
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

//########################################################################
//######################################################################
//#			UINT322ASCII				       #
//######################################################################
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

//########################################################################
//######################################################################
//#			UINT642ASCII				       #
//######################################################################
// uint642ascii takes 8 bytes and returns a string with a numeric
// decimal interpretation of the found data as ascii-string
func uint642ascii(payload []byte) string {
	if len(payload) != 8 {
		return "Err in CAN-Frame, data must be 8 bytes."
	}
	data := binary.LittleEndian.Uint64(payload)
	return strconv.FormatUint(uint64(data), 10)
}

func ascii2uint64(payload string) []byte {
	tmp, _ := strconv.Atoi(payload)
	number := uint64(tmp)
	a := make([]byte, 8)
	binary.LittleEndian.PutUint64(a, number)
	return a
}

//########################################################################
//######################################################################
//#             bytecolor2colorcode
//######################################################################
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
