package main

import (
	"bufio"        // Reader
	"encoding/csv" // CSV Management
	"flag"
	"github.com/c3re/can2mqtt/convertfunctions"
	"io"  // EOF const
	"log" // error management
	"log/slog"
	"os"      // open files
	"strconv" // parse strings
	"sync"
)

var pairFromID map[uint32]*can2mqtt    // c2m pair (lookup from ID)
var pairFromTopic map[string]*can2mqtt // c2m pair (lookup from Topic)
var debugLog bool
var canInterface, mqttConnection, configFile string
var dirMode = 0 // directional modes: 0=bidirectional 1=can2mqtt only 2=mqtt2can only [-d]
var wg sync.WaitGroup

func main() {
	log.SetFlags(0)

	flag.BoolVar(&debugLog, "v", false, "show (very) verbose debug log")
	flag.StringVar(&canInterface, "c", "can0", "which socket-can interface to use")
	flag.StringVar(&mqttConnection, "m", "tcp://localhost:1883", "which mqtt-broker to use. Example: tcp://user:password@broker.hivemq.com:1883")
	flag.StringVar(&configFile, "f", "can2mqtt.csv", "which config file to use")
	flag.IntVar(&dirMode, "d", 0, "direction mode\n0: bidirectional (default)\n1: can2mqtt only\n2: mqtt2can only")
	flag.Parse()

	if dirMode < 0 || dirMode > 2 {
		slog.Error("got invalid value for -d. Valid values are 0 (bidirectional), 1 (can2mqtt only) or 2 (mqtt2can only)", "d", dirMode)
	}

	if debugLog {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.Info("Starting can2mqtt", "mqtt-config", mqttConnection, "can-interface", canInterface, "can2mqtt.csv", configFile, "dir-mode", dirMode, "debug", debugLog)
	wg.Add(1)
	go canStart(canInterface) // epic parallel shit ;-)
	mqttStart(mqttConnection)
	readC2MPFromFile(configFile)
	wg.Wait()
}

// this functions opens, parses and extracts information out
// of the can2mqtt.csv
func readC2MPFromFile(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		slog.Error("can2mqtt.csv could not be opened", "filename", filename, "error", err)
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
			slog.Warn("skipping line", "filename", filename, "error", err)
			continue
		}
		line, _ := r.FieldPos(0)
		tmp, err := strconv.ParseUint(record[0], 10, 32)
		if err != nil {
			slog.Warn("skipping line, malformed can-ID", "error", err, "line", line)
			continue
		}
		canID := uint32(tmp)
		convMode := record[1]
		topic := record[2]
		if isIDInSlice(canID) {
			slog.Warn("skipping line, duplicate ID", "id", canID, "line", line)
			continue
		}
		if isTopicInSlice(topic) {
			slog.Warn("skipping line duplicate topic", "topic", topic, "line", line)
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
			// Int methods come here now
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
		slog.Debug("extracted pair", "id", c2mp.canId, "convertmode", c2mp.convMethod, "topic", c2mp.mqttTopic)
	}
}

func isIDInSlice(canId uint32) bool {
	return pairFromID[canId] != nil
}

func isTopicInSlice(mqttTopic string) bool {
	return pairFromTopic[mqttTopic] != nil
}

// get the corresponding topic for an ID
func getTopicFromId(canId uint32) string {
	return pairFromID[canId].mqttTopic
}
