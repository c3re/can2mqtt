package can2mqtt

import (
	"fmt"
	CAN "github.com/brendoncarroll/go-can"
	"strconv"
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
	if convertMethod == "bytes2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode ascii2bytes (reverse of %s)\n", convertMethod)
		}
		data, len = ascii2bytes(payload)
	} else if convertMethod == "byte2dec" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode dec2byte (reverse of %s)\n", convertMethod)
		}
		data[0] = dec2byte(payload)
		len = 1
	} else if convertMethod == "openorclosed2oneorzero" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmodus dec2byte (reverse of %s)\n", convertMethod)
		}
		data, len = oneorzero2openorclosed(payload)
	} else {
		if dbg {
			fmt.Printf("convertfunctions: convertmode %s not found. using fallback ascii2byte\n", convertMethod)
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
func convert2MQTT(id int, length int, payload [8]byte) string {
	convertMethod := getConvId(id)
	if convertMethod == "bytes2ascii" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode bytes2asciii\n")
		}

		return bytes2ascii(uint32(length), payload)
	} else if convertMethod == "byte2dec" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode byte2dec\n")
		}
		return byte2dec(payload[0])
	} else if convertMethod == "openorclosed2oneorzero" {
		if dbg {
			fmt.Printf("convertfunctions: using convertmode byte2dec\n")
		}
		return openorclosed2oneorzero(payload[0])
	} else {
		if dbg {
			fmt.Printf("convertfunctions: convertmode %s not found. using fallback bytes2ascii\n", convertMethod)
		}
		return bytes2ascii(uint32(length), payload)
	}
}

// some convertmodes
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

// byte2dec takes exactly one byte and returns a string with a
// numeric decimal interpretation of the found data
func byte2dec(payload byte) string {
	return strconv.FormatInt(int64(payload), 10)
}

func dec2byte(payload string) byte {
	tmp, err := strconv.ParseInt(payload, 8, 10)
	if err != nil {
		return byte(255)
	}
	return byte(tmp)
}

// openorclosed2oneorzero. takes one byte with a "o" for the return of "1"
// or everything else for a return of "0"
// 1 -> opened
// 0 and everything else -> closed
func openorclosed2oneorzero(payload byte) string {
	if string(payload) == "o" {
		return "1"
	} else {
		return "0"
	}
}

func oneorzero2openorclosed(payload string) ([8]byte, uint32) {
	if payload == "1" {
		return ascii2bytes("open   ")
	} else {
		return ascii2bytes("closed ")
	}
}
