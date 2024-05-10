// Package internal of c3re/can2mqtt contains some tools for bridging a CAN-Interface
// and a mqtt-network
package internal

import (
	"github.com/brutella/can"
	"log/slog"
	"os"
	"sync"
)

var csi []uint32       // subscribed IDs slice
var csiLock sync.Mutex // CAN subscribed IDs Mutex
var bus *can.Bus       // CAN-Bus pointer

// initializes the CANBus Interface and enters an infinite
// loop that reads CAN-frames after that.
func canStart(canInterface string) {

	var err error
	slog.Debug("canbus: initializing CAN-Bus", "interface", canInterface)
	bus, err = can.NewBusForInterfaceWithName(canInterface)
	if err != nil {
		slog.Error("canbus: error while initializing CAN-Bus", "error", err)
		os.Exit(1)
	}
	slog.Info("canbus: connected to CAN")
	slog.Debug("canbus: registering handler")
	bus.SubscribeFunc(handleCANFrame)
	slog.Debug("canbus: starting receive loop")
	// will not return if everything is fine
	err = bus.ConnectAndPublish()
	if err != nil {
		slog.Error("canbus: error while processing CAN-Bus", "error", err)
		os.Exit(1)
	}
}

func handleCANFrame(frame can.Frame) {
	frame.ID &= 0x1FFFFFFF // discard flags, we are only interested in the ID
	var idSub = false      // indicates, whether the id was subscribed or not
	for _, i := range csi {
		if i == frame.ID {
			slog.Debug("canbus: received subscribed frame", "id", frame.ID)
			go handleCAN(frame)
			idSub = true
			break
		}
	}
	if !idSub {
		slog.Debug("canbus: ignored unsubscribed frame", "id", frame.ID)
	}
}

// Unsubscribe a CAN-ID
func canSubscribe(id uint32) {
	csiLock.Lock()
	csi = append(csi, id)
	csiLock.Unlock()
	slog.Debug("canbus: successfully subscribed CAN-ID", "id", id)
}

// expects a CANFrame and sends it
func canPublish(frame can.Frame) {
	slog.Debug("canbus: sending CAN-Frame", "frame", frame)
	// Check if ID is using more than 11-Bits:
	if frame.ID >= 0x800 {
		// if so, enable extended frame format
		frame.ID |= 0x80000000
	}
	err := bus.Publish(frame)
	if err != nil {
		slog.Error("canbus: error while publishing CAN-Frame", "error", err)
	}
}
