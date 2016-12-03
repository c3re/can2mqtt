package main

import (
	"bufio"                                // Reader
	"encoding/csv"                         // CSV Management
	"fmt"                                  // printfoo
	CAN "github.com/brendoncarroll/go-can" // CAN-Bus Binding
	"io"                                   // EOF const
	"log"                                  // error management
	"os"                                   // open files
	"strconv"                              // parse strings
	"sync"
)

type can2mqtt struct {
	canId      int
	convMethod string
	mqttTopic  string
}

var can2mqttPairs []can2mqtt
var dbg bool = false
var wg sync.WaitGroup

func main() {
	if len(os.Args) == 1 {
		printHelp()
	} else {
		if len(os.Args) == 5 {
			if os.Args[4] == "-v" {
				dbg = true
			}
		}
		if os.Args[1] == "test-mqtt" {
			dbg = true
			if dbg {
				fmt.Printf("main: Starting MQTT-Test:\n")
			}
			MQTTStart(os.Args[2])
			MQTTPublish("test", os.Args[3])
		} else if os.Args[1] == "test-can" {
			dbg = true
			if dbg {
				fmt.Printf("main: Starting CAN-Bus-Test:\n")
			}
			CANStart(os.Args[2])
			data, datalength := ascii2bytes(os.Args[3])
			cf := CAN.CANFrame{ID: 112, Len: datalength, Data: data}
			CANPublish(cf)
		} else {
			wg.Add(1)
			go CANStart(os.Args[2]) // epic parallel shit ;-)
			MQTTStart(os.Args[3])
			readC2MPFromFile(os.Args[1])
			wg.Wait()
		}
	}
}

func printHelp() {
	fmt.Printf("Willkommen zur CAN2MQTT Bridge!\n\n")
	fmt.Printf("Usage: can2mqtt <file> <CAN-Interface> <MQTT-Connect>\n")
	fmt.Printf("<file>: entweder eine Datei oder test-can oder test-mqtt\n")
	fmt.Printf("<CAN-Interface>: Ein CAN-Interface z.B. can0\n")
	fmt.Printf("<MQTT-Connect>: Connectring fuer MQTT. Beispiel: tcp://localhost:1883\n")
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
		can2mqttPairs = append(can2mqttPairs, can2mqtt{i, record[1], record[2]})
		MQTTSubscribe(record[2])
		CANSubscribe(uint32(i))
	}
	if dbg {
		fmt.Printf("main: Die folgenden CAN-MQTT Kombinationen wurden gelesen:\n")
		fmt.Printf("main: CAN-ID\t\tConversion Mode\t\tMQTT-Topic")
		for _, c2mp := range can2mqttPairs {
			fmt.Printf("main: %d\t\t%s\t\t%s\n", c2mp.canId, c2mp.convMethod, c2mp.mqttTopic)
		}
	}
}

func IsInSlice(canId int, mqttTopic string) bool {
	for _, c2mp := range can2mqttPairs {
		if c2mp.canId == canId || c2mp.mqttTopic == mqttTopic {
			if dbg {
				fmt.Printf("main: Die ID %d oder das Topic %s wurden bereits angegeben!\n", canId, mqttTopic)
			}
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
