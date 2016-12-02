package main

import (
	"bufio"		// Reader
	"encoding/csv"	// CSV Management
	"fmt"		// printfoo
	"log"		// error management
	"os"		// open files
	"io"		// EOF const
	"strconv"	// parse strings
	CAN "github.com/brendoncarroll/go-can" // CAN-Bus Binding
	"time"		// sleepy stuff
)

type can2mqtt struct {
	canId int
	convMethod string
	mqttTopic string
}
var can2mqttPairs []can2mqtt

func main() {
	if len(os.Args) == 1 {
		printHelp()
	} else {
		if os.Args[1] == "test-mqtt" {
			fmt.Println("main: Starting MQTT-Test:")
			MQTTStart(os.Args[2])
			MQTTPublish("test", os.Args[3])
		} else if os.Args[1] == "test-can" {
			fmt.Println("main: Starting CAN-Bus-Test:")
			CANStart(os.Args[2])
			data, datalength := ascii2byte(os.Args[3])
			cf := CAN.CANFrame{ID: 112, Len: datalength, Data: data}
			CANPublish(cf)
		} else {
			go CANStart(os.Args[2]) // epic parallel shit ;-)
			MQTTStart(os.Args[3])
			readC2MPFromFile(os.Args[1])
			time.Sleep(-1)
			time.Sleep(1000000 * time.Millisecond)
		}
	}
}

func printHelp() {
	fmt.Println("Willkommen zur CAN2MQTT Bridge!")
	fmt.Println()
	fmt.Println("Usage: can2mqtt <file> <CAN-Interface> <MQTT-Connect>")
	fmt.Println("<file>: entweder eine Datei oder test-can oder test-mqtt")
	fmt.Println("<CAN-Interface>: Ein CAN-Interface z.B. can0")
	fmt.Println("<MQTT-Connect>: Connectring fuer MQTT. Beispiel: tcp://localhost:1883")
}

func readC2MPFromFile(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(bufio.NewReader(file))
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		i, err := strconv.Atoi(record[0])
		if IsInSlice(i, record[2]) {
			panic("main: Jede ID und jedes Topic darf maximal einmal auftreten!")
		}
		can2mqttPairs = append(can2mqttPairs, can2mqtt{i,record[1],record[2]})
		MQTTSubscribe(record[2])
		CANSubscribe(uint32(i))
	}
	fmt.Println("main: Die folgenden CAN-MQTT Kombinationen wurden gelesen:")
	fmt.Println("main: CAN-ID\t\tConversion Mode\t\tMQTT-Topic")
	for _, c2mp := range can2mqttPairs {
		fmt.Printf("main: %d\t\t%s\t\t%s\n", c2mp.canId, c2mp.convMethod, c2mp.mqttTopic)
	}
}

func IsInSlice(canId int, mqttTopic string) bool {
	for _, c2mp := range can2mqttPairs {
		if c2mp.canId == canId || c2mp.mqttTopic == mqttTopic {
			fmt.Printf("main: Die ID %d oder das Topic %s wurden bereits angegeben!\n", canId, mqttTopic)
			return true
		}
	}
	return false
}

func getTopic(canId int) string {
	for _, c2mp := range can2mqttPairs {
		if c2mp.canId == canId {
			return c2mp.mqttTopic
		}
	}
	// Fehlerfall
	return "-1"
}

func getConvTopic(topic string) string {
	for _, c2mp := range can2mqttPairs {
		if c2mp.mqttTopic == topic {
			return c2mp.convMethod
		}
	}
	// Fehlerfall
	return "-1"
}


func getId(mqttTopic string) int {
	for _, c2mp := range can2mqttPairs {
		if c2mp.mqttTopic == mqttTopic {
			return c2mp.canId
		}
	}
	// Fehlerfall
	return -1
}


func getConvId(canId int) string {
	for _, c2mp := range can2mqttPairs {
		if c2mp.canId == canId {
			return c2mp.convMethod
		}
	}
	// Fehlerfall
	return "-1"
}

