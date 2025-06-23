package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/brutella/can"
	"github.com/jaster-prj/can2mqtt/config"
)

type CanOptions struct {
	Interface string
}

type CanListener struct {
	running     bool
	bus         *can.Bus // CAN-Bus pointer
	options     *CanOptions
	lock        sync.Mutex          // CAN subscribed IDs Mutex
	csi         map[uint32]struct{} // subscribed IDs slice
	publishCan  chan can.Frame
	publishMqtt chan MqttPublish
	routing     *config.Routing
	converter   *ConverterFactory
}

func NewCanListener() *CanListener {
	return &CanListener{
		running:   false,
		bus:       nil,
		lock:      sync.Mutex{},
		csi:       map[uint32]struct{}{},
		routing:   nil,
		converter: nil,
	}
}

func (c *CanListener) WithOptions(options *CanOptions) *CanListener {
	c.options = options
	return c
}

func (c *CanListener) WithRouting(routing *config.Routing) *CanListener {
	c.routing = routing
	return c
}

func (c *CanListener) WithConverter(converter *ConverterFactory) *CanListener {
	c.converter = converter
	return c
}

func (c *CanListener) WithPublishCanChannel(publishCan chan can.Frame) *CanListener {
	c.publishCan = publishCan
	return c
}

func (c *CanListener) WithPublishMqttChannel(publishMqtt chan MqttPublish) *CanListener {
	c.publishMqtt = publishMqtt
	return c
}

// initializes the CANBus Interface and enters an infinite
// loop that reads CAN-frames after that.
func (c *CanListener) Run() {

	var err error
	if c.options == nil {
		slog.Error("canbus: error while initializing CAN-Bus: no options defined")
		os.Exit(1)
	}
	slog.Debug("canbus: initializing CAN-Bus", "interface", c.options.Interface)
	c.bus, err = can.NewBusForInterfaceWithName(c.options.Interface)
	if err != nil {
		slog.Error("canbus: error while initializing CAN-Bus", "interface", c.options.Interface, "error", err)
		os.Exit(1)
	}
	slog.Info("canbus: connected to CAN")
	slog.Debug("canbus: registering handler")
	c.bus.SubscribeFunc(c.handleCANFrame)
	slog.Debug("canbus: starting receive loop")
	// will not return if everything is fine
	go func() {
		c.running = true
		err = c.bus.ConnectAndPublish()
		if err != nil {
			slog.Error("canbus: error while processing CAN-Bus", "error", err)
			os.Exit(1)
		}
		c.running = false
		close(c.publishCan)
	}()

	go func() {
		for {
			select {
			case frame := <-c.publishCan:
				slog.Debug("canbus: sending CAN-Frame", "frame", frame)
				// Check if ID is using more than 11-Bits:
				if frame.ID >= 0x800 {
					// if so, enable extended frame format
					frame.ID |= 0x80000000
				}
				err := c.bus.Publish(frame)
				if err != nil {
					slog.Error("canbus: error while publishing CAN-Frame", "error", err)
				}
			}
			if c.publishCan == nil {
				break
			}
		}
	}()
}

func (c *CanListener) Stop() {
	if c.running == true {
		c.bus.Disconnect()
		close(c.publishCan)
		c.running = false
	}
}

func (c *CanListener) UpdateConfiguration(addRoutes []config.Route, delRoutes []config.Route) {
	slog.Info("CanListener UpdateConfiguration")
	for _, route := range delRoutes {
		c.Unsubscribe(route.CanID)
	}
	for _, route := range addRoutes {
		c.Subscribe(route.CanID)
	}
}

// subscribe to a new ident
func (c *CanListener) Subscribe(ident string) {
	id, err := strconv.Atoi(ident)
	if err != nil {
		return
	}
	c.lock.Lock()
	c.csi[uint32(id)] = struct{}{}
	c.lock.Unlock()
	slog.Debug("canbus: successfully subscribed CAN-ID", "id", id)
}

// unsubscribe a ident
func (c *CanListener) Unsubscribe(ident string) {
	id, err := strconv.Atoi(ident)
	if err != nil {
		return
	}
	c.lock.Lock()
	if _, ok := c.csi[uint32(id)]; ok {
		delete(c.csi, uint32(id))
	}
	c.lock.Unlock()
	slog.Debug("canbus: successfully unsubscribed CAN-ID", "id", id)
}

func (c *CanListener) handleCANFrame(frame can.Frame) {
	frame.ID &= 0x1FFFFFFF // discard flags, we are only interested in the ID
	var idSub = false      // indicates, whether the id was subscribed or not
	if _, ok := c.csi[frame.ID]; ok {
		slog.Debug("canbus: received subscribed frame", "id", frame.ID)
		go c.handleCAN(frame)
		idSub = true
	}
	if !idSub {
		slog.Debug("canbus: ignored unsubscribed frame", "id", frame.ID)
	}
}

// handleCAN is the standard receive handler for CANFrames
// and does the following:
// 1. calling standard convert function: convert2MQTT
// 2. sending the message
func (c *CanListener) handleCAN(cf can.Frame) {
	slog.Debug("received CANFrame", "id", cf.ID, "len", cf.Length, "data", cf.Data)

	route, err := c.routing.GetRouteByCanId(fmt.Sprintf("%d", cf.ID))
	if err != nil {
		slog.Warn("GetRouteByCanId error", "error", err)
		return
	}
	if route.Direction == config.MQTT2CAN {
		return
	}
	var converterStr string
	if route.Converter != nil {
		converterStr = *route.Converter
	} else {
		converterStr = "none"
	}
	converter, err := c.converter.GetConverter(converterStr)
	if err != nil {
		slog.Warn("GetConverter error", "error", err)
		return
	}
	mqttPayload, err := converter.ToMqtt(cf)
	if err != nil {
		slog.Warn("conversion to MQTT message unsuccessful", "convertmode", converter, "error", err)
		return
	}
	topic := route.Topic
	c.publishMqtt <- MqttPublish{
		Topic:   topic,
		Payload: mqttPayload,
	}
	// this is the most common log-message, craft with care...
	// slog.Debug("CAN -> MQTT", "ID", cf.ID, "len", cf.Length, "data", cf.Data, "convertmode", converter, "topic", topic, "message", mqttPayload)
}
