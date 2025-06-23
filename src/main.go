package main

import (
	// Reader
	// CSV Management
	"crypto/tls"
	"fmt"

	// EOF const
	// error management
	"log/slog" // open files
	// parse strings
	"sync"

	"github.com/brutella/can"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/jaster-prj/can2mqtt/config"
)

var (
	debugLog bool
	version  = "dev"
	wg       sync.WaitGroup
)

func onConnectionLost(c MQTT.Client, e error) {
	slog.Debug("MQTT: Connection lost.", e)
}

func onReconnecting(c MQTT.Client, o *MQTT.ClientOptions) {
	slog.Debug("MQTT: Reconnecting")
}

func main() {
	appConfig := config.GetConfiguration()

	if appConfig.LogLevel != nil {
		slog.SetLogLoggerLevel(*appConfig.LogLevel)
	}

	canPublishChannel := make(chan can.Frame, 10)
	mqttPublishChannel := make(chan MqttPublish, 10)

	slog.Info("Starting can2mqtt", "version", version, "mqtt-config", appConfig.MqttConnection, "can-interface", appConfig.Device, "debug", debugLog)
	wg.Add(1)
	converterFactory := NewConverterFactory()
	routing := config.NewRouting()
	updater := config.NewUpdate().
		WithRouting(routing)
	canListener := NewCanListener().
		WithConverter(converterFactory).
		WithOptions(&CanOptions{Interface: appConfig.Device}).
		WithPublishCanChannel(canPublishChannel).
		WithPublishMqttChannel(mqttPublishChannel).
		WithRouting(routing)

	protocol := "tcp"
	if appConfig.MqttConnection.Protocol != nil {
		protocol = *appConfig.MqttConnection.Protocol
	}
	port := ":1883"
	if appConfig.MqttConnection.Port != nil {
		port = fmt.Sprintf(":%d", *appConfig.MqttConnection.Port)
	}

	// create MQTT Client
	clientSettings := MQTT.NewClientOptions().AddBroker(protocol + "://" + appConfig.MqttConnection.Url + port).SetResumeSubs(true)
	clientSettings.SetClientID(appConfig.MqttConnection.ClientName)
	if appConfig.MqttConnection.Username != nil {
		clientSettings.SetUsername(*appConfig.MqttConnection.Username)
	}
	if appConfig.MqttConnection.Password != nil {
		clientSettings.SetPassword(*appConfig.MqttConnection.Password)
	}
	clientSettings.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	mqttListener := NewMqttListener().
		WithOptions(clientSettings).
		WithPublishCanChannel(canPublishChannel).
		WithPublishMqttChannel(mqttPublishChannel).
		WithRouting(routing).
		WithConverter(converterFactory)

	updater.RegisterCallback(mqttListener.UpdateConfiguration)
	updater.RegisterCallback(canListener.UpdateConfiguration)

	canListener.Run()
	mqttListener.Run()

	mqttListener.SubscribeConfig(updater.ConfigUpdate)
	wg.Wait()
}
