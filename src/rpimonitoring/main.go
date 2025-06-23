package main

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MqttListener struct {
	running bool
	client  *MQTT.Client
	options *MQTT.ClientOptions
}

func NewMqttListener() *MqttListener {
	return &MqttListener{
		running: false,
		client:  nil,
		options: nil,
	}
}

func getPiTemperatureFromSys() (float64, error) {
	data, err := os.ReadFile("/sys/class/thermal/thermal_zone0/temp")
	if err != nil {
		return 0, err
	}
	temperatureStr := strings.TrimSpace(string(data))
	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err != nil {
		return 0, err
	}
	return temperature / 1000, nil
}

func main() {
	appConfig := GetConfiguration()

	if appConfig.LogLevel != nil {
		slog.SetLogLoggerLevel(*appConfig.LogLevel)
	}

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

	client := MQTT.NewClient(clientSettings)
	slog.Debug("mqtt: starting connection", "connectString", clientSettings.Servers)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("mqtt: could not connect to mqtt", "error", token.Error())
		os.Exit(1)
	}
	slog.Info("mqtt: connected to mqtt")

	for {
		temp, err := getPiTemperatureFromSys()
		if err != nil {
			continue
		}
		token := client.Publish(fmt.Sprintf("/%s/temperature", appConfig.Device), 0, false, fmt.Sprintf("%.2f°C\n", temp))
		token.Wait()
		time.Sleep(120 * time.Second)
	}
}
