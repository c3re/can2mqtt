package internal

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
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
		if !found {
			fmt.Println("Whoops, there is an issue with your MQTT-connectString:")
			fmt.Println("suppliedString: ", suppliedString)
			fmt.Println("userPasswordHost: ", userPasswordHost)
		}
		user, pw, found = strings.Cut(userPassword, ":")
		if !found {
			fmt.Println("Whoops, there is an issue with your MQTT-connectString:")
			fmt.Println("suppliedString: ", suppliedString)
			fmt.Println("userPasswordHost: ", userPasswordHost)
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
	if dbg {
		fmt.Printf("mqtthandler: starting connection to: %s\n", connectString)
	}
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("mqttHandler: Oh no an error occurred...")
		panic(token.Error())
	}
	if dbg {
		fmt.Printf("mqttHandler: connection established!\n")
	}
}

// credentialsProvider
func userPwCredProv() (username, password string) {
	return user, pw
}

// subscribe to a new topic
func mqttSubscribe(topic string) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("mqtthandler: error while subscribing: %s\n", topic)
	}
	if dbg {
		fmt.Printf("mqtthandler: successfully subscribed: %s\n", topic)
	}
}

// unsubscribe a topic
func mqttUnsubscribe(topic string) {
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Printf("mqtthandler: Error while unsuscribing :%s\n", topic)
	}
	if dbg {
		fmt.Printf("mqtthandler: successfully unsubscribed %s\n", topic)
	}
}

// publish a new message
func mqttPublish(topic string, payload string) {
	if dbg {
		fmt.Printf("mqtthandler: sending message: \"%s\" to topic: \"%s\"\n", payload, topic)
	}
	mqttUnsubscribe(topic)
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	if dbg {
		fmt.Printf("mqtthandler: message was transmitted successfully!.\n")
	}
	mqttSubscribe(topic)
}
