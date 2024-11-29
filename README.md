# can2mqtt
can2mqtt based on https://github.com/c3re/can2mqtt. It connects a CAN-Bus to an MQTT broker and vice versa. You create pairs of CAN-Bus IDs and MQTT topics and choose a conversion mode between them.

## Installation
The latest compiled binary is available in the [releases](https://github.com/jaster-prj/can2mqtt/releases/latest).
The binaries are statically linked and have no further dependencies. On a Raspberry you can install via:
```
wget https://github.com/jaster-prj/can2mqtt/releases/download/v2.3.0/can2mqtt-v2.3.0-linux-arm -O can2mqtt
chmod +x can2mqtt
./can2mqtt
```

## Usage

can2mqtt has to be configured to connect to can-device and mqtt broker. This configuration can be loaded from File or environmental variables
```
{
    "loglevel": "debug",
    "device": "/dev/can0",
    "mqttconnection": {
        "protocol": "wss",       //default tcp
        "url": "broker.url.com", 
        "port": 443,             //default 1883
        "username": "<username>",
        "password": "<password>"
    }
}
```
```
export CONFIG_FILE="config.json"
export LOGLEVEL="info"
export DEVICE="/dev/can0"
export MQTTURL="broker.url.com"
export MQTTPORT="443"
export MQTTUSERNAME="<username>"
export MQTTPASSWORD="<password>"
```
 
To start service execute binary
```
./can2mqtt
```

## can2mqtt gateway config
Connected to the mqtt broker, can2mqtt will subscribe to topic <br>/gateway/routes</br>. This topic expects json-formated list of routes. To be shure, can2mqtt will receive latest configuration, all configuration should be published in retain mode.
can2mqtt gets newest version of routes-configuration and will compare changes with latest configuration. All changed routes will be unsubcribed and new ones subscribed. 
Direction can be set by direction-value. 0 means BIDIRECTIONAL, 1 publishes only mqtt changes to can, 2 publishes only can messages to mqtt

Example:
```
[
  {
    "canid": "112",
    "topic": "/can/msg/112",
    "direction": 0,
    "converter": "none"
  },
  {
    "canid": "113",
    "topic": "/can/msg/113",
    "direction": 1,
    "converter": "uint322ascii"
  },
  {
    "canid": "114",
    "topic": "/can/msg/114",
    "direction": 2,
    "converter": "pixelbin2ascii"
  },
]
```

Explanation for the 1st Element: Can-Message with CAN-ID 112 (decimal) will be published to /can/msg/112, when message will be received. Due to direction is bidirectional, changes from /can/msg/112 will also be published on can-bus with CAN-ID 112. The none-converter means, received data will not be modified.

## convert-modes
Here they are:

| convertmode           | description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
|-----------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `none`                | does not convert anything. It just takes a bunch of bytes and hands it over to the other side. If you want to send strings, this will be your choice. If you have a mqtt payload that is longer than eight bytes, only the first eight bytes will be send via CAN.                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| `bytecolor2colorcode` | Converts an bytearray of 3 bytes to hexadecimal colorcode                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `pixelbin2ascii`      | This mode was designed to address colorized pixels. MQTT-wise you can insert a string like "<0-255> #RRGGBB" which will be converted to 4 byte on the CAN-BUS the first byte will be the number of the LED 0-255 and bytes 1, 2, 3 are the color of red, green and blue.                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `16bool2ascii`        | Interprets two bytes can-wise and publishes them as 16 boolean values to mqtt                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| `{i}[u]int{b}2ascii`  | `i` is the amount of instances of numbers in the CAN-Frame/the MQTT-message. Valid instance amounts are currently 1,2,4 and 8. Although other combinations are possible. Might be added in the future, if the need arises. `b` is the size of each number in bits. Supported values are 8,16,32 and 64. You can use either one unsigned or signed integers. This are all possible combinations:  `int82ascii`, `2int82ascii`, `4int82ascii`, `8int82ascii`, `int162ascii`, `2int162ascii`, `4int162ascii`, `int322ascii`, `2int322ascii`, `int642ascii`, `uint82ascii`, `2uint82ascii`, `4uint82ascii`, `8uint82ascii`, `uint162ascii`, `2uint162ascii`, `4uint162ascii`, `uint322ascii`, `2uint322ascii`, `uint642ascii` |


## How do I get a socket-can interface?
You can either use a hardware interface or setup a virtual socket-can interface.
### Hardware interface
There are many articles on the internet on how to get a can interface in Linux. The important part here is that you get a socket-can interface in the end. But in most cases this is possible. For example the MCP2515 chip can be used with the SPI interface of a raspberry to create a socket-can interface. There is a driver for the ELM327 chip too. Serial-to-CAN converters can be used too via `slcand`.
### Virtual interface
For testing purposes or to just get you going you can use `vcan` a virtual can interface that comes with Linux itself. You can configure it for example like this:
```bash
sudo ip link add dev vcan0 type vcan
sudo ip link set vcan0 up
```
Now you can use your new interface `vcan0` as socket-can interface with can2mqtt.

## Debugging
To debug the behaviour of can2mqtt you need to be able to send and receive CAN frames and MQTT messages. For MQTT I recommend [mosquitto](https://mosquitto.org/) with its `mosquitto_pub` and `mosquitto_sub` commands. For CAN i recommend [can-utils]() with its tools `cansend` and `candump`.

## Add a convert-Mode
If you want to add a convert-Mode think about a name. This is the name that you can later refer to when you want to
use your convert-Mode in the broker config json. Now, use the file `src/convertmode/mymode.go` as a template for your own convertmode. Copy that file to `src/convertmode/<yournewmode>.go`. Now change all occurrences of "MyMode" with your preferred Name (Lets say `YourNewMode` in this example). Note that it has to start with an upper-case letter, so that it is usable outside of this package. Next you have to write three functions (implement the `ConvertMode` interface):
1. A conversion method from CAN -> MQTT: `ToMqtt(input can.Frame) ([]byte, error)`
2. A conversion method from MQTT -> CAN: `ToCan(input []byte) (can.Frame, error)`
3. A `String() string` method that reports the name of that convertmode. This method is used in some log-messages

Your almost done, the last step is to "register" your new convertmode. To do so add the following line to [`src/converter.go#L13`](./src/converter.go#L13)
```go
convertModeFromString[convertmode.<YourNewMode>{}.String()] = convertmode.<YourNewMode>{}
```

Now you can use your new convertmode in your broker config. Use the string that your return in the `String()` function as the name of the convertmode. In the `mymode.go` code this is `"mymode"`.

Good luck & happy hacking ✌
