package internal

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log/slog"
	"os"
	"strings"
)

var client MQTT.Client
var user, pw string

// uses the connectString to establish a connection to the MQTT
// broker
func mqttStart(suppliedString string) {
	connectString := suppliedString
	if strings.Contains(suppliedString, "@") {
		// looks like authentication is required for this server
		userPasswordHost := strings.TrimPrefix(suppliedString, "tcp://")
		userPassword, host, found := strings.Cut(userPasswordHost, "@")
		user, pw, found = strings.Cut(userPassword, ":")
		if !found {
			slog.Error("mqtthandler: missing colon(:) between username and password", "connect string", suppliedString)
			os.Exit(1)
		}
		connectString = "tcp://" + host
	}
	clientSettings := MQTT.NewClientOptions().AddBroker(connectString)
	clientSettings.SetClientID("CAN2MQTT")
	clientSettings.SetDefaultPublishHandler(handleMQTT)
	if strings.Contains(suppliedString, "@") {
		clientSettings.SetCredentialsProvider(userPwCredProv)
	}
	client = MQTT.NewClient(clientSettings)
	slog.Debug("mqtthandler: starting connection", "connectString", connectString)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("mqtthandler: could not connect to mqtt", "error", token.Error())
		os.Exit(1)
	}
	slog.Info("mqtthandler: connected to mqtt")
}

// credentialsProvider
func userPwCredProv() (username, password string) {
	return user, pw
}

// subscribe to a new topic
func mqttSubscribe(topic string) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		slog.Error("mqtthandler: error subscribing", "error", token.Error())
	}
	slog.Debug("mqtthandler: subscribed", "topic", topic)
}

// unsubscribe a topic
func mqttUnsubscribe(topic string) {
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		slog.Error("mqtthandler: error unsubscribing", "error", token.Error())
	}
	slog.Debug("mqtthandler: unsubscribed", "topic", topic)
}

// publish a new message
func mqttPublish(topic string, payload []byte) {
	mqttUnsubscribe(topic)
	slog.Debug("mqtthandler: publishing message", "payload", payload, "topic", topic)
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	slog.Debug("mqtthandler: published message", "payload", payload, "topic", topic)
	mqttSubscribe(topic)
}
