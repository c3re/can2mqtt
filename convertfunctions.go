package main

import (
	CAN "github.com/brendoncarroll/go-can"
	"fmt"
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
		fmt.Printf("convertfunctions: Benutze convertmodus ascii2byte (reverse of %s)\n", convertMethod)
		data, len = ascii2byte(payload)
	} else {
		fmt.Printf("convertfunctions: convertmodus %s nicht gefunden. Benutze Fallback ascii2byte\n", convertMethod)
		data, len = ascii2byte(payload)
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
	if convertMethod == "ascii2byte" {
		fmt.Printf("convertfunctions: Benutze convertmodus byte2asciii\n")

		return byte2ascii(uint32(length), payload)
	} else {
		fmt.Printf("convertfunctions: convertmodus %s nicht gefunden. Benutze Fallback byte2ascii\n", convertMethod)
		return byte2ascii(uint32(length), payload)
	}
}

// Eine der Konvertermethoden
func ascii2byte(payload string) ([8]byte, uint32) {
	var returner [8]byte
	var i uint32 = 0
	for ; int(i) < len(payload) && i < 8; i++ {
		returner[i] = payload[i]
	}
	return returner, i
}

func byte2ascii(length uint32, payload [8]byte) string {
	return string(payload[:length])
}

