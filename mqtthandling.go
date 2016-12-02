package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var client MQTT.Client

func MQTTStart(connectString string) {
	clientsettings := MQTT.NewClientOptions().AddBroker(connectString)
	clientsettings.SetClientID("CAN2MQTT")
	clientsettings.SetDefaultPublishHandler(handleMQTT)
	client = MQTT.NewClient(clientsettings)
	fmt.Printf("mqtthandler: Starte Verbindung nach: %s\n", connectString)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("mqtthandler: Oh no an error occured...")
		panic(token.Error())
	}
	fmt.Printf("mqtthandler: Verbindung erfolgreich hergestellt\n")
}

func MQTTSubscribe(topic string) {
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("mqtthandler: Error beim abonnieren von %s\n", topic, token.Error())
	}
	fmt.Printf("mqtthandler: Topic: %s erfolgreich abonniert\n", topic)
}

func MQTTUnsubscribe(topic string) {
	if token := client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Printf("mqtthandler: Error beim unsuscriben von %s\n", topic, token.Error())
	}
	fmt.Printf("mqtthandler: Topic: %s erfolgreich unsubscribed\n", topic)
}

func MQTTPublish(topic string, payload string) {
	fmt.Printf("mqtthandler: Sende Nachricht: \"%s\" an Topic: \"%s\"\n", payload, topic)
	MQTTUnsubscribe(topic)
	token := client.Publish(topic, 0, false, payload)
	token.Wait()
	fmt.Println("mqtthandler: Nachricht erfolgreich gesendet.")
	MQTTSubscribe(topic)
}
