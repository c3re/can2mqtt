package internal

import (
	"bufio"        // Reader
	"encoding/csv" // CSV Management
	"fmt"          // print :)
	"io"           // EOF const
	"log"          // error management
	"os"           // open files
	"strconv"      // parse strings
	"log/slog"
	"github.com/brutella/can"
	"github.com/c3re/can2mqtt/internal/convertfunctions"
	"sync"
)

type convertToCan func(input []byte) (can.Frame, error)
type convertToMqtt func(input can.Frame) ([]byte, error)

type ConvertMode interface {
	convertToCan
	convertToMqtt
}

// can2mqtt is a struct that represents the internal type of
// one line of the can2mqtt.csv file. It has
// the same three fields as the can2mqtt.csv file: CAN-ID,
// conversion method and MQTT-Topic.
type can2mqtt struct {
	canId      uint32
	convMethod string
	toCan      convertToCan
	toMqtt     convertToMqtt
	mqttTopic  string
}

var pairFromID map[uint32]*can2mqtt    // c2m pair (lookup from ID)
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
	slog.SetLogLoggerLevel(slog.LevelDebug)
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
	log.SetFlags(0)
	slog.Info("Starting can2mqtt", "mqtt-config", cs, "can-interface", ci, "can2mqtt.csv", c2mf, "dir-mode", dirMode, "debug", dbg)
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
		slog.Error("main: can2mqtt.csv could not be opened", "filename", filename, "error", err)
		os.Exit(1)
	}

	r := csv.NewReader(bufio.NewReader(file))
	r.FieldsPerRecord = 3
	pairFromID = make(map[uint32]*can2mqtt)
	pairFromTopic = make(map[string]*can2mqtt)
	for {
		record, err := r.Read()
		// Stop at EOF.
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Warn("main: skipping line", "filename", filename, "error", err)
			continue
		}
		line, _ := r.FieldPos(0)
		tmp, err := strconv.ParseUint(record[0], 10, 32)
		if err != nil {
			slog.Warn("main: skipping line, malformed can-ID", "error", err, "line", line)
			continue
		}
		canID := uint32(tmp)
		convMode := record[1]
		topic := record[2]
		if isIDInSlice(canID) {
			slog.Warn("main: skipping line, duplicate ID", "id", canID, "line", line)
			continue
		}
		if isTopicInSlice(topic) {
			slog.Warn("main: skipping line duplicate topic", "topic", topic, "line", line)
			continue
		}
		switch convMode {
		case "16bool2ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.SixteenBool2AsciiToCan,
				toMqtt:     convertfunctions.SixteenBool2AsciiToMqtt,
			}
		case "uint82ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.Uint82AsciiToCan,
				toMqtt:     convertfunctions.Uint82AsciiToMqtt,
			}
		case "uint162ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.Uint162AsciiToCan,
				toMqtt:     convertfunctions.Uint162AsciiToMqtt,
			}
		case "uint322ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.Uint322AsciiToCan,
				toMqtt:     convertfunctions.Uint322AsciiToMqtt,
			}
		case "uint642ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.Uint642AsciiToCan,
				toMqtt:     convertfunctions.Uint642AsciiToMqtt,
			}
		case "2uint322ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.TwoUint322AsciiToCan,
				toMqtt:     convertfunctions.TwoUint322AsciiToMqtt,
			}
		case "4uint162ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.FourUint162AsciiToCan,
				toMqtt:     convertfunctions.FourUint162AsciiToMqtt,
			}
		case "4uint82ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.FourUint82AsciiToCan,
				toMqtt:     convertfunctions.FourUint82AsciiToMqtt,
			}
		case "8uint82ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.EightUint82AsciiToCan,
				toMqtt:     convertfunctions.EightUint82AsciiToMqtt,
			}
			// Int methodes come here now
		case "int82ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.Int82AsciiToCan,
				toMqtt:     convertfunctions.Int82AsciiToMqtt,
			}
		case "int162ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.Int162AsciiToCan,
				toMqtt:     convertfunctions.Int162AsciiToMqtt,
			}
		case "int322ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.Int322AsciiToCan,
				toMqtt:     convertfunctions.Int322AsciiToMqtt,
			}
		case "int642ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.Int642AsciiToCan,
				toMqtt:     convertfunctions.Int642AsciiToMqtt,
			}
		case "2int322ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.TwoInt322AsciiToCan,
				toMqtt:     convertfunctions.TwoInt322AsciiToMqtt,
			}
		case "4int162ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.FourInt162AsciiToCan,
				toMqtt:     convertfunctions.FourInt162AsciiToMqtt,
			}
		case "4int82ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.FourInt82AsciiToCan,
				toMqtt:     convertfunctions.FourInt82AsciiToMqtt,
			}
		case "8int82ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.EightInt82AsciiToCan,
				toMqtt:     convertfunctions.EightInt82AsciiToMqtt,
			}
		case "bytecolor2colorcode":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.ByteColor2ColorCodeToCan,
				toMqtt:     convertfunctions.ByteColor2ColorCodeToMqtt,
			}
		case "pixelbin2ascii":
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.PixelBin2AsciiToCan,
				toMqtt:     convertfunctions.PixelBin2AsciiToMqtt,
			}
		default:
			pairFromID[canID] = &can2mqtt{
				canId:      canID,
				convMethod: convMode,
				mqttTopic:  topic,
				toCan:      convertfunctions.NoneToCan,
				toMqtt:     convertfunctions.NoneToMqtt,
			}
		}
		pairFromTopic[topic] = pairFromID[canID]
		mqttSubscribe(topic) // TODO move to append function
		canSubscribe(canID)  // TODO move to append function
	}

	for _, c2mp := range pairFromID {
		slog.Debug("main: extracted pair", "id", c2mp.canId, "convertmode", c2mp.convMethod, "topic", c2mp.mqttTopic)
	}
}

func isIDInSlice(canId uint32) bool {
	return pairFromID[canId] != nil
}

func isTopicInSlice(mqttTopic string) bool {
	return pairFromTopic[mqttTopic] != nil
}

// get the corresponding ID for a given topic
func getIdFromTopic(topic string) uint32 {
	return pairFromTopic[topic].canId
}

// get the conversion mode for a given topic
func getConvModeFromTopic(topic string) string {
	return pairFromTopic[topic].convMethod
}

// get the convertMode for a given ID
func getConvModeFromId(canId uint32) string {
	return pairFromID[canId].convMethod
}

// get the corresponding topic for an ID
func getTopicFromId(canId uint32) string {
	return pairFromID[canId].mqttTopic
}
