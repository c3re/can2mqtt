// minimal main package, only commandline argument parsing and a printHelp() function
package main

import (
	"fmt" // print
	C2M "github.com/c3re/can2mqtt/internal"
	"os" // args
)

// Parses commandline arguments
func main() {
	conf := true
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-v":
			C2M.SetDbg(true)
		case "-c":
			i++
			C2M.SetCi(os.Args[i])
		case "-m":
			i++
			C2M.SetCs(os.Args[i])
		case "-f":
			i++
			C2M.SetC2mf(os.Args[i])
		case "-d":
			i++
			C2M.SetConfDirMode(os.Args[i])
		default:
			i = len(os.Args)
			conf = false
			printHelp()
		}
	}
	if conf {
		C2M.Start()
	}
}

// help function (obvious...)
func printHelp() {
	_, _ = fmt.Fprintf(os.Stderr, "welcome to the CAN2MQTT bridge!\n\n")
	_, _ = fmt.Fprintf(os.Stderr, "Usage: can2mqtt [-f <file>] [-c <CAN-Interface>] [-m <MQTT-Connect>] [-v] [-h] [-d <dirMode>]\n")
	_, _ = fmt.Fprintf(os.Stderr, "<file>: a can2mqtt.csv file\n")
	_, _ = fmt.Fprintf(os.Stderr, "<CAN-Interface>: a CAN-Interface e.g. can0\n")
	_, _ = fmt.Fprintf(os.Stderr, "<MQTT-Connect>: connectstring for MQTT. e.g.: tcp://[user:pass@]localhost:1883\n")
	_, _ = fmt.Fprintf(os.Stderr, "<dirMode>: directional Mode 0 = bidirectional, 1 = can2mqtt only, 2 = mqtt2can only\n")
}
