package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"net/url"
	"os"
)

var client MQTT.Client
var user, pw string

// uses the connectString to establish a connection to the MQTT
// broker
func mqttStart(URL string) {
	// parse the supplied URL
	u, err := url.Parse(URL)
	if err != nil {
		slog.Error("while parsing URL", "url", URL, "error", err)
		os.Exit(1)
	}

	// create MQTT Client
	clientSettings := MQTT.NewClientOptions().AddBroker(u.Scheme + "://" + u.Host)
	clientSettings.SetClientID("can2mqtt")
	clientSettings.SetDefaultPublishHandler(handleMQTT)
	if u.User != nil {
		clientSettings.SetUsername(u.User.Username())
		password, passwdSet := u.User.Password()
		if passwdSet {
			clientSettings.SetPassword(password)
		}
	}

	client = MQTT.NewClient(clientSettings)
	slog.Debug("mqtt: starting connection", "connectString", URL)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("mqtt: could not connect to mqtt", "error", token.Error())
		os.Exit(1)
	}
	slog.Info("mqtt: connected to mqtt")
}

// subscribe to a new topic
func mqttSubscribe(topic string) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		slog.Error("mqtt: error subscribing", "error", token.Error())
	}
	slog.Debug("mqtt: subscribed", "topic", topic)
}

// unsubscribe a topic
func mqttUnsubscribe(topic string) {
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		slog.Error("mqtt: error unsubscribing", "error", token.Error())
	}
	slog.Debug("mqtt: unsubscribed", "topic", topic)
}

// publish a new message
func mqttPublish(topic string, payload []byte) {
	mqttUnsubscribe(topic)
	slog.Debug("mqtt: publishing message", "payload", payload, "topic", topic)
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	slog.Debug("mqtt: published message", "payload", payload, "topic", topic)
	mqttSubscribe(topic)
}