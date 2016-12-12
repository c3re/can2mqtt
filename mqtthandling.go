package can2mqtt

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var client MQTT.Client

// uses the connectString to establish a connection to the MQTT
// broker
func mqttStart(connectString string) {
	clientsettings := MQTT.NewClientOptions().AddBroker(connectString)
	clientsettings.SetClientID("CAN2MQTT")
	clientsettings.SetDefaultPublishHandler(handleMQTT)
	client = MQTT.NewClient(clientsettings)
	if dbg {
		fmt.Printf("mqtthandler: starting connection to: %s\n", connectString)
	}
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("mqtthandler: Oh no an error occured...")
		panic(token.Error())
	}
	if dbg {
		fmt.Printf("mqtthandler: connection established!\n")
	}
}

// subscribe to a new topic
func mqttSubscribe(topic string) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("mqtthandler: error while subscribing: %s\n", topic, token.Error())
	}
	if dbg {
		fmt.Printf("mqtthandler: successfully subscribed: %s\n", topic)
	}
}

// unsubscribe a topic
func mqttUnsubscribe(topic string) {
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Printf("mqtthandler: Error while unsuscribing :%s\n", topic, token.Error())
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
