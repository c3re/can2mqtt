package main

import (
	"fmt"
	CAN "github.com/brendoncarroll/go-can"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	)
// Standard Receivehandler fuer CAN-Frames. Diese Funktion macht:
// 1. standardconverterfunktion aufrufen convert2MQTT
// 2. Message wegschicken
func handleCAN(cf CAN.CANFrame) {
	if dbg { fmt.Printf("receivehandler: Habe CAN Frame empfangen: ID: %d, Laenge: %d, Inhalt %s\n", cf.ID, cf.Len, cf.Data) }
	mqttPayload := convert2MQTT(int(cf.ID), int(cf.Len), cf.Data)
	if dbg { fmt.Printf("receivehandler: Konvertierter String lautet: %s\n", mqttPayload) }
	topic := getTopic(int(cf.ID))
	MQTTPublish(topic, mqttPayload)
	fmt.Printf("ID: %d Length: %d Data: %X -> Topic: \"%s\" Message: \"%s\"\n", cf.ID, cf.Len, cf.Data, topic, mqttPayload)
}

// Standard Receivehandler fuer MQTT-Messages. Diese Funktion macht:
// 1. standardconverterfunktion aufrufen convert2CAN
// 2. Message wegschicken
func handleMQTT(cl MQTT.Client, msg MQTT.Message) {
	if dbg { fmt.Printf("receivehandler: Habe MQTT-Message empfangen: Topic: %s, Msg: %s\n", msg.Topic(), msg.Payload()) }
	cf := convert2CAN(msg.Topic(), string(msg.Payload()))
	CANPublish(cf)
	fmt.Printf("ID: %d Length: %d Data: %X <- Topic: \"%s\" Message: \"%s\"\n", cf.ID, cf.Len, cf.Data, msg.Topic(), msg.Payload())
}
