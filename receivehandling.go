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
	fmt.Printf("receivehandler: Habe CAN Frame empfangen: ID: %d, Laenge: %d, Inhalt %s\n", cf.ID, cf.Len, cf.Data)
	mqttPayload := convert2MQTT(int(cf.ID), int(cf.Len), cf.Data)
	fmt.Printf("receivehandler: Konvertierter String lautet: %s\n", mqttPayload)
	MQTTPublish(getTopic(int(cf.ID)), mqttPayload)
}

// Standard Receivehandler fuer MQTT-Messages. Diese Funktion macht:
// 1. standardconverterfunktion aufrufen convert2CAN
// 2. Message wegschicken
func handleMQTT(cl MQTT.Client, msg MQTT.Message) {
	fmt.Printf("receivehandler: Habe MQTT-Message empfangen: Topic: %s, Msg: %s\n", msg.Topic(), msg.Payload())
	cf := convert2CAN(msg.Topic(), string(msg.Payload()))
	CANPublish(cf)
}
