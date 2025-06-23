package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/brutella/can"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jaster-prj/can2mqtt/config"
)

const (
	GATEWAY_ROUTES = "/gateway/routes"
)

type MqttPublish struct {
	Topic   string
	Payload []byte
}

type MqttListener struct {
	running     bool
	client      *MQTT.Client
	options     *MQTT.ClientOptions
	publishMqtt chan MqttPublish
	publishCan  chan can.Frame
	routing     *config.Routing
	converter   *ConverterFactory
}

func NewMqttListener() *MqttListener {
	return &MqttListener{
		running:   false,
		client:    nil,
		options:   nil,
		routing:   nil,
		converter: nil,
	}
}

func (m *MqttListener) WithOptions(options *MQTT.ClientOptions) *MqttListener {
	m.options = options
	return m
}

func (m *MqttListener) WithRouting(routing *config.Routing) *MqttListener {
	m.routing = routing
	return m
}

func (m *MqttListener) WithConverter(converter *ConverterFactory) *MqttListener {
	m.converter = converter
	return m
}

func (m *MqttListener) WithPublishCanChannel(publishCan chan can.Frame) *MqttListener {
	m.publishCan = publishCan
	return m
}

func (m *MqttListener) WithPublishMqttChannel(publishMqtt chan MqttPublish) *MqttListener {
	m.publishMqtt = publishMqtt
	return m
}

// uses the connectString to establish a connection to the MQTT
// broker
func (m *MqttListener) Run() {

	options := m.options
	options.SetDefaultPublishHandler(m.handleMQTT)
	client := MQTT.NewClient(m.options)
	slog.Debug("mqtt: starting connection", "connectString", m.options.Servers)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("mqtt: could not connect to mqtt", "error", token.Error())
		os.Exit(1)
	}
	slog.Info("mqtt: connected to mqtt")
	m.running = true
	m.client = &client

	go func() {
		for {
			select {
			case publish := <-m.publishMqtt:
				m.Unsubscribe(publish.Topic)
				slog.Debug("mqtt: publishing message", "payload", publish.Payload, "topic", publish.Topic)
				token := client.Publish(publish.Topic, 0, false, publish.Payload)
				token.Wait()
				slog.Debug("mqtt: published message", "payload", publish.Payload, "topic", publish.Topic)
				m.Subscribe(publish.Topic)
			}
			if m.publishMqtt == nil {
				break
			}
		}
	}()
}

func (m *MqttListener) Stop() {
	if m.running == true {
		if m.client != nil {
			(*m.client).Disconnect(0)
		}
		m.running = false
	}
}

func (m *MqttListener) UpdateConfiguration(addRoutes []config.Route, delRoutes []config.Route) {
	slog.Info("MqttListener UpdateConfiguration")
	for _, route := range delRoutes {
		m.Unsubscribe(route.Topic)
	}
	for _, route := range addRoutes {
		m.Subscribe(route.Topic)
	}
}

// subscribe to a new ident
func (m *MqttListener) SubscribeConfig(callback func(config []byte)) {
	if m.client == nil {
		slog.Debug("client not connected")
	}
	client := (*m.client)
	if token := client.Subscribe(
		GATEWAY_ROUTES,
		1,
		func(_ MQTT.Client, msg MQTT.Message) {
			callback(msg.Payload())
		},
	); token.Wait() && token.Error() != nil {
		slog.Error("mqtt: error subscribing", "error", token.Error())
	}
	slog.Debug("mqtt: subscribed", "topic", GATEWAY_ROUTES)
}

// subscribe to a new ident
func (m *MqttListener) Subscribe(ident string) {
	if m.client == nil {
		slog.Debug("client not connected")
	}
	client := (*m.client)
	if token := client.Subscribe(ident, 0, nil); token.Wait() && token.Error() != nil {
		slog.Error("mqtt: error subscribing", "error", token.Error())
	}
	slog.Debug("mqtt: subscribed", "topic", ident)
}

// unsubscribe a ident
func (m *MqttListener) Unsubscribe(ident string) {
	if m.client == nil {
		slog.Debug("client not connected")
	}
	client := (*m.client)
	if token := client.Unsubscribe(ident); token.Wait() && token.Error() != nil {
		slog.Error("mqtt: error unsubscribing", "error", token.Error())
	}
	slog.Debug("mqtt: unsubscribed", "topic", ident)
}

// handleMQTT is the standard receive handler for MQTT
// messages and does the following:
// 1. calling the standard convert function: convert2CAN
// 2. sending the message
func (m *MqttListener) handleMQTT(_ MQTT.Client, msg MQTT.Message) {
	slog.Debug("received message", "topic", msg.Topic(), "payload", msg.Payload())

	route, err := m.routing.GetRouteByMqttTopic(msg.Topic())
	if route.Direction == config.CAN2MQTT {
		return
	}
	var converterStr string
	if route.Converter != nil {
		converterStr = *route.Converter
	} else {
		converterStr = "none"
	}
	converter, err := m.converter.GetConverter(converterStr)
	if err != nil {
		slog.Warn("GetConverter error", "error", err)
		return
	}
	cf, err := converter.ToCan(msg.Payload())
	if err != nil {
		slog.Warn("conversion to CAN-Frame unsuccessful", "convertmode", converter, "error", err)
		return
	}
	canID, err := strconv.Atoi(route.CanID)
	if err != nil {
		slog.Warn("conversion to CAN-ID failed", "convertmode", converter, "error", err)
		return
	}
	cf.ID = uint32(canID)
	m.publishCan <- cf
	// slog.Debug("CAN <- MQTT", "ID", cf.ID, "len", cf.Length, "data", cf.Data, "convertmode", converter, "topic", msg.Topic(), "message", msg.Payload())
}
