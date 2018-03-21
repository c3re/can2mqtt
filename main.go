package can2mqtt

import (
	"bufio"        // Reader
	"encoding/csv" // CSV Management
	"fmt"          // printfoo
	"io"           // EOF const
	"log"          // error management
	"os"           // open files
	"strconv"      // parse strings
	"sync"
)

// can2mqtt is a struct that represents the internal type of
// one line of the can2mqtt.csv file. Therefore it has
// the same three fields as the can2mqtt.csv file: CAN-ID,
// conversion method and MQTT-Topic.
type can2mqtt struct {
	canId      uint32
	convMethod string
	mqttTopic  string
}

var can2mqttPairs []can2mqtt           // representation of can2mqtt.csv
var dbg bool = false                   // verbose on of [-v]
var ci string = "can0"                 // the CAN-Interface [-c]
var cs string = "tcp://localhost:1883" // mqtt-connectstring [-m]
var c2mf string = "can2mqtt.csv"       // path to the can2mqtt.csv [-f]
var conf bool = true                   // represents wether a running conf
// is set or not
var wg sync.WaitGroup

// SetDbg decides wether there is really verbose output or
// just standard information output. Default is false.
func SetDbg(v bool) {
	dbg = v
}

// SetCi sets the CAN-Interface to use for the CAN side
// of the bridge. Default is: can0.
func SetCi(c string) {
	ci = c
}

// SetC2mf expects a string which is a path to a can2mqtt.csv file
// Default is: can2mqtt.csv
func SetC2mf(f string) {
	c2mf = f
}

// SetCs sets the MQTT connectstring which contains: protocol,
// hostname and port. Default is: tcp://localhost:1883
func SetCs(s string) {
	cs = s
}

// Start is the function that should be called after debug-level
// connectstring, can interface and can2mqtt file have been set.
// Start takes care of everything that happens after that.
// It starts the CAN-Bus connection and the MQTT-Connection. It
// parses the can2mqtt.csv file and from there everything takes
// its course...
func Start() {
	wg.Add(1)
	go canStart(ci) // epic parallel shit ;-)
	mqttStart(cs)
	readC2MPFromFile(c2mf)
	wg.Wait()
}

// this functions opens, parses and extracts information out
// of the can2mqtt.csv
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
		u, err := strconv.ParseUint(record[0], 10, 64)
		if isInSlice(uint32(u), record[2]) {
			panic("main: each ID and each topic is only allowed once!")
		}
		can2mqttPairs = append(can2mqttPairs, can2mqtt{uint32(u), record[1], record[2]})
		mqttSubscribe(record[2])
		canSubscribe(uint32(u))
	}
	if dbg {
		fmt.Printf("main: the following CAN-MQTT pairs have been extracted:\n")
		fmt.Printf("main: CAN-ID\t\t conversion mode\t\tMQTT-topic")
		for _, c2mp := range can2mqttPairs {
			fmt.Printf("main: %d\t\t%s\t\t%s\n", c2mp.canId, c2mp.convMethod, c2mp.mqttTopic)
		}
	}
}

// check function to check if a topic or an ID is in the slice
func isInSlice(canId uint32, mqttTopic string) bool {
	for _, c2mp := range can2mqttPairs {
		if c2mp.canId == canId || c2mp.mqttTopic == mqttTopic {
			if dbg {
				fmt.Printf("main: The ID %d or the Topic %s is already in the list!\n", canId, mqttTopic)
			}
			return true
		}
	}
	return false
}

// get the corresponding topic for an ID
func getTopic(canId uint32) string {
	for _, c2mp := range can2mqttPairs {
		if c2mp.canId == canId {
			return c2mp.mqttTopic
		}
	}
	// Fehlerfall
	return "-1"
}

// get the conversion mode for a given topic
func getConvTopic(topic string) string {
	for _, c2mp := range can2mqttPairs {
		if c2mp.mqttTopic == topic {
			return c2mp.convMethod
		}
	}
	// Fehlerfall
	return "-1"
}

// get the correspondig ID for a given topic
func getId(mqttTopic string) uint32 {
	for _, c2mp := range can2mqttPairs {
		if c2mp.mqttTopic == mqttTopic {
			return c2mp.canId
		}
	}
	// Fehlerfall
	return 0
}

// get the convertode for a given ID
func getConvId(canId uint32) string {
	for _, c2mp := range can2mqttPairs {
		if c2mp.canId == canId {
			return c2mp.convMethod
		}
	}
	// Fehlerfall
	return "-1"
}
