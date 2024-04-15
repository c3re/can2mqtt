// Package can2mqtt contains some tools for bridging a CAN-Interface
// and a mqtt-network
package internal

import (
	"fmt"
	"github.com/brutella/can"
	"log"
	"sync"
)

var csi []uint32       // subscribed IDs slice
var csiLock sync.Mutex // CAN subscribed IDs Mutex
var bus *can.Bus       // CAN-Bus pointer

// initializes the CANBus Interface and enters an infinite
// loop that reads CAN-frames after that.
func canStart(canInterface string) {

	var err error
	if dbg {
		fmt.Printf("canbushandler: initializing CAN-Bus interface %s\n", canInterface)
	}
	bus, err = can.NewBusForInterfaceWithName(canInterface)
	if err != nil {
		if dbg {
			fmt.Printf("canbushandler: error while activating CAN-Bus interface: %s\n", canInterface)
		}
		log.Fatal(err)
	}
	bus.SubscribeFunc(handleCANFrame)
	err = bus.ConnectAndPublish()
	if err != nil {
		if dbg {
			fmt.Printf("canbushandler: error while activating CAN-Bus interface: %s\n", canInterface)
		}
		log.Fatal(err)
	}
}

func handleCANFrame(frame can.Frame) {
	frame.ID &= 0x1FFFFFFF // discard flags, we are only interested in the ID
	var idSub = false      // indicates, whether the id was subscribed or not
	for _, i := range csi {
		if i == frame.ID {
			if dbg {
				fmt.Printf("canbushandler: ID %d is in subscribed list, calling receivehadler.\n", frame.ID)
			}
			go handleCAN(frame)
			idSub = true
			break
		}
	}
	if !idSub {
		if dbg {
			fmt.Printf("canbushandler: ID:%d was not subscribed. /dev/nulled that frame...\n", frame.ID)
		}
	}
}

// Unsubscribe a CAN-ID
func canSubscribe(id uint32) {
	csiLock.Lock()
	csi = append(csi, id)
	csiLock.Unlock()
	if dbg {
		fmt.Printf("canbushandler: mutex lock+unlock successful. subscribed to ID:%d\n", id)
	}
}

// expects a CANFrame and sends it
func canPublish(frame can.Frame) {
	if dbg {
		fmt.Println("canbushandler: sending CAN-Frame: ", frame)
	}
	// Check if ID is using more than 11-Bits:
	if frame.ID >= 0x800 {
		// if so, enable extended frame format
		frame.ID |= 0x80000000
	}
	err := bus.Publish(frame)
	if err != nil {
		if dbg {
			fmt.Printf("canbushandler: error while transmitting the CAN-Frame.\n")
		}
		log.Fatal(err)
	}
}
