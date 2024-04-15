package internal

import (
	"bufio"        // Reader
	"encoding/csv" // CSV Management
	"fmt"          // print :)
	"io"           // EOF const
	"log"          // error management
	"os"           // open files
	"strconv"      // parse strings
	"sync"
)

// can2mqtt is a struct that represents the internal type of
// one line of the can2mqtt.csv file. It has
// the same three fields as the can2mqtt.csv file: CAN-ID,
// conversion method and MQTT-Topic.
type can2mqtt struct {
	canId      int
	convMethod string
	mqttTopic  string
}

var pairFromID map[int]*can2mqtt       // c2m pair (lookup from ID)
var pairFromTopic map[string]*can2mqtt // c2m pair (lookup from Topic)
var dbg = false                        // verbose on off [-v]
var ci = "can0"                        // the CAN-Interface [-c]
var cs = "tcp://localhost:1883"        // mqtt-connect-string [-m]
var c2mf = "can2mqtt.csv"              // path to the can2mqtt.csv [-f]
var dirMode = 0                        // directional modes: 0=bidirectional 1=can2mqtt only 2=mqtt2can only [-d]
var wg sync.WaitGroup

// SetDbg decides whether there is really verbose output or
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

// SetCs sets the MQTT connect-string which contains: protocol,
// hostname and port. Default is: tcp://localhost:1883
func SetCs(s string) {
	cs = s
}

// SetConfDirMode sets the dirMode
func SetConfDirMode(s string) {
	if s == "0" {
		dirMode = 0
	} else if s == "1" {
		dirMode = 1
	} else if s == "2" {
		dirMode = 2
	} else {
		_ = fmt.Errorf("error: got invalid value for -d (%s). Valid values are 0 (bidirectional), 1 (can2mqtt only) or 2 (mqtt2can only)", s)
	}
}

// Start is the function that should be called after debug-level
// connect-string, can interface and can2mqtt file have been set.
// Start takes care of everything that happens after that.
// It starts the CAN-Bus connection and the MQTT-Connection. It
// parses the can2mqtt.csv file and from there everything takes
// its course...
func Start() {
	fmt.Println("Starting can2mqtt")
	fmt.Println()
	fmt.Println("MQTT-Config:  ", cs)
	fmt.Println("CAN-Config:   ", ci)
	fmt.Println("can2mqtt.csv: ", c2mf)
	fmt.Print("dirMode:       ", dirMode, " (")
	if dirMode == 0 {
		fmt.Println("bidirectional)")
	}
	if dirMode == 1 {
		fmt.Println("can2mqtt only)")
	}
	if dirMode == 2 {
		fmt.Println("mqtt2can only)")
	}
	fmt.Print("Debug-Mode:    ")
	if dbg {
		fmt.Println("yes")
	} else {
		fmt.Println("no")
	}
	fmt.Println()
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
	pairFromID = make(map[int]*can2mqtt)
	pairFromTopic = make(map[string]*can2mqtt)
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		canID, err := strconv.Atoi(record[0])
		convMode := record[1]
		topic := record[2]
		if isInSlice(canID, topic) {
			panic("main: each ID and each topic is only allowed once!")
		}
		pairFromID[canID] = &can2mqtt{
			canId:      canID,
			convMethod: convMode,
			mqttTopic:  topic,
		}
		pairFromTopic[topic] = pairFromID[canID]
		mqttSubscribe(topic)        // TODO move to append function
		canSubscribe(uint32(canID)) // TODO move to append function
	}
	if dbg {
		fmt.Printf("main: the following CAN-MQTT pairs have been extracted:\n")
		fmt.Printf("main: CAN-ID\t\t conversion mode\t\tMQTT-topic\n")
		for _, c2mp := range pairFromID {
			fmt.Printf("main: %d\t\t%s\t\t%s\n", c2mp.canId, c2mp.convMethod, c2mp.mqttTopic)
		}
	}
}

// check function to check if a topic or an ID is in the slice
func isInSlice(canId int, mqttTopic string) bool {
	if pairFromID[canId] != nil {
		if dbg {
			fmt.Printf("main: The ID %d or the Topic %s is already in the list!\n", canId, mqttTopic)
		}
		return true
	}
	if pairFromTopic[mqttTopic] != nil {
		if dbg {
			fmt.Printf("main: The ID %d or the Topic %s is already in the list!\n", canId, mqttTopic)
		}
		return true
	}
	return false
}

// get the corresponding ID for a given topic
func getIdFromTopic(topic string) int {
	return pairFromTopic[topic].canId
}

// get the conversion mode for a given topic
func getConvModeFromTopic(topic string) string {
	return pairFromTopic[topic].convMethod
}

// get the convertMode for a given ID
func getConvModeFromId(canId int) string {
	return pairFromID[canId].convMethod
}

// get the corresponding topic for an ID
func getTopicFromId(canId int) string {
	return pairFromID[canId].mqttTopic
}
