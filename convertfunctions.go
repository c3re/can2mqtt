package main

import (
	CAN "github.com/brendoncarroll/go-can"
	"fmt"
	"strconv"
	)


// Standardfunktion macht folgendes:
// 1. topic und payload entgegennehmen
// 2. zugehoerige converterfunktion herraussuchen
// 3. zugehoerige ID herraussuchen
// 3. Konvertierung durchfuehren 
// 4. CANFrame bauen
// 5. zurueckgeben
func convert2CAN (topic, payload string) CAN.CANFrame {
	convertMethod := getConvTopic(topic)
	var Id uint32 = uint32(getId(topic))
	var data [8]byte
	var len uint32
	if convertMethod == "bytes2ascii" {
		if dbg { fmt.Printf("convertfunctions: Benutze convertmodus ascii2bytes (reverse of %s)\n", convertMethod) }
		data, len = ascii2bytes(payload)
	} else if convertMethod == "byte2dec" {
		if dbg { fmt.Printf("convertfunctions: Benutze convertmodus dec2byte (reverse of %s)\n", convertMethod) }
		data[0] = dec2byte(payload)
		len = 1
	} else {
		if dbg { fmt.Printf("convertfunctions: convertmodus %s nicht gefunden. Benutze Fallback ascii2byte\n", convertMethod) }
		data, len = ascii2bytes(payload)
	}
	mycf := CAN.CANFrame{ID: Id, Len: len, Data: data}
	return mycf
}

// Standardfunktion macht folgendes:
// 1. id und payload entgegennehmen
// 2. zugehoerige converterfunktion herraussuchen
// 3. Konvertierung durchfuehren 
// 4. string bauen
// 5. zurueckgeben
func convert2MQTT (id int, length int, payload [8]byte) string {
	convertMethod := getConvId(id)
	if convertMethod == "bytes2ascii" {
		if dbg { fmt.Printf("convertfunctions: Benutze convertmodus bytes2asciii\n") }

		return bytes2ascii(uint32(length), payload)
	} else if convertMethod == "byte2dec" {
		if dbg { fmt.Printf("convertfunctions: Benutze convertmodus byte2dec\n") }
		return byte2dec(payload[0])
	} else {
		if dbg { fmt.Printf("convertfunctions: convertmodus %s nicht gefunden. Benutze Fallback bytes2ascii\n", convertMethod) }
		return bytes2ascii(uint32(length), payload)
	}
}

// Eine der Konvertermethoden
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

// byte2dec. Nimmt ein byte entgegen und gibt einen String mit dezimaler Interpretation aus
func byte2dec(payload byte) string {
	return strconv.FormatInt(int64(payload), 10)
}

func dec2byte(payload string) byte{
	tmp, err := strconv.ParseInt(payload, 8, 10)
	if err != nil {
		return byte(255)
	}
	return byte(tmp)
}
