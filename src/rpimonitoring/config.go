package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/jaster-prj/can2mqtt/common"
)

type AppConfig struct {
	LogLevel       *slog.Level `json:"loglevel,omitempty"`
	Device         string      `json:"device"`
	MqttConnection MqttConfig  `json:"mqttconnection"`
}

type MqttConfig struct {
	ClientName string  `json:"clientname"`
	Protocol   *string `json:"protocol,omitempty"`
	Url        string  `json:"url"`
	Port       *int    `json:"port,omitempty"`
	Username   *string `json:"username,omitempty"`
	Password   *string `json:"password,omitempty"`
}

func GetConfiguration() AppConfig {
	configFile := "config.json"
	if os.Getenv("CONFIG_FILE") != "" {
		configFile = os.Getenv("CONFIG_FILE")
	}
	appConfig := AppConfig{}
	if _, err := os.Stat(configFile); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Parse Config from destination: %s\n", configFile)
		appConfig = parseConfigFile(configFile)
	}
	if os.Getenv("LOGLEVEL") != "" {
		appConfig.LogLevel.UnmarshalText([]byte(os.Getenv("LOGLEVEL")))
	}
	if os.Getenv("DEVICE") != "" {
		appConfig.Device = os.Getenv("DEVICE")
	}
	if os.Getenv("MQTTURL") != "" {
		appConfig.MqttConnection.Url = os.Getenv("MQTTURL")
	}
	if os.Getenv("MQTTPORT") != "" {
		port, err := strconv.Atoi(os.Getenv("MQTTPORT"))
		if err == nil {
			appConfig.MqttConnection.Port = &port
		}
	}
	if os.Getenv("MQTTUSERNAME") != "" {
		appConfig.MqttConnection.Username = common.POINTER(os.Getenv("MQTTUSERNAME"))
	}
	if os.Getenv("MQTTPASSWORD") != "" {
		appConfig.MqttConnection.Password = common.POINTER(os.Getenv("MQTTPASSWORD"))
	}
	return appConfig
}

func parseConfigFile(configFile string) AppConfig {
	var appConfig AppConfig
	jsonData, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Println(err)
		return AppConfig{}
	}
	err = json.Unmarshal(jsonData, &appConfig)
	if err != nil {
		fmt.Println(err)
		return AppConfig{}
	}
	return appConfig
}
