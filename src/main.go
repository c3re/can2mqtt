package main

import (
	"bufio"        // Reader
	"encoding/csv" // CSV Management
	"flag"
	"github.com/c3re/can2mqtt/convertmode"
	"io"  // EOF const
	"log" // error management
	"log/slog"
	"os"      // open files
	"strconv" // parse strings
	"sync"
)

var (
	pairFromID                               map[uint32]*can2mqtt // c2m pair (lookup from ID)
	pairFromTopic                            map[string]*can2mqtt // c2m pair (lookup from Topic)
	convertModeFromString                    map[string]ConvertMode
	debugLog                                 bool
	canInterface, mqttConnection, configFile string
	version                                  = "dev"
	dirMode                                  = BIDIRECTIONAL // directional modes: 0=bidirectional 1=can2mqtt only 2=mqtt2can only [-d]
	wg                                       sync.WaitGroup
)

func main() {
	log.SetFlags(0)

	flag.BoolVar(&debugLog, "v", false, "show (very) verbose debug log")
	flag.StringVar(&canInterface, "c", "can0", "which socket-can interface to use")
	flag.StringVar(&mqttConnection, "m", "tcp://localhost:1883", "which mqtt-broker to use. Example: tcp://user:password@broker.hivemq.com:1883")
	flag.StringVar(&configFile, "f", "can2mqtt.csv", "which config file to use")
	flag.IntVar(&dirMode, "d", 0, "direction mode\n0: bidirectional (default)\n1: can2mqtt only\n2: mqtt2can only")
	flag.Parse()

	if dirMode < BIDIRECTIONAL || dirMode > MQTT2CAN_ONLY {
		slog.Error("got invalid value for -d. Valid values are 0 (bidirectional), 1 (can2mqtt only) or 2 (mqtt2can only)", "d", dirMode)
	}

	if debugLog {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	slog.Info("Starting can2mqtt", "version", version, "mqtt-config", mqttConnection, "can-interface", canInterface, "can2mqtt.csv", configFile, "dir-mode", dirMode, "debug", debugLog)
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
	convertModeFromString = make(map[string]ConvertMode)

	// initialize all convertModes
	convertModeFromString[convertmode.None{}.String()] = convertmode.None{}
	convertModeFromString[convertmode.SixteenBool2Ascii{}.String()] = convertmode.SixteenBool2Ascii{}
	convertModeFromString[convertmode.PixelBin2Ascii{}.String()] = convertmode.PixelBin2Ascii{}
	convertModeFromString[convertmode.ByteColor2ColorCode{}.String()] = convertmode.ByteColor2ColorCode{}
	convertModeFromString[convertmode.MyMode{}.String()] = convertmode.MyMode{}
	// Dynamically create int and uint convertmodes
	for _, bits := range []uint{8, 16, 32, 64} {
		for _, instances := range []uint{1, 2, 4, 8} {
			if bits*instances <= 64 {
				// int
				cmi, _ := convertmode.NewInt2Ascii(instances, bits)
				convertModeFromString[cmi.String()] = cmi
				// uint
				cmu, _ := convertmode.NewUint2Ascii(instances, bits)
				convertModeFromString[cmu.String()] = cmu
			}
		}
	}
	if debugLog {
		for _, cm := range convertModeFromString {
			slog.Debug("convertmode initialized", "convertmode", cm)
		}
	}
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
		if pairFromID[canID] != nil {
			slog.Warn("skipping line, duplicate ID", "id", canID, "line", line)
			continue
		}
		if pairFromTopic[topic] != nil {
			slog.Warn("skipping line duplicate topic", "topic", topic, "line", line)
			continue
		}

		if convertModeFromString[convMode] == nil {
			slog.Warn("skipping line, unsupported convertMode ", "convertMode", convMode, "line", line)
			continue
		}

		pairFromID[canID] = &can2mqtt{
			canId:       canID,
			convertMode: convertModeFromString[convMode],
			mqttTopic:   topic,
		}
		pairFromTopic[topic] = pairFromID[canID]
		mqttSubscribe(topic) // TODO move to append function
		canSubscribe(canID)  // TODO move to append function
	}

	for _, c2mp := range pairFromID {
		slog.Debug("extracted pair", "id", c2mp.canId, "convertmode", c2mp.convertMode, "topic", c2mp.mqttTopic)
	}
}