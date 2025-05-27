# can2mqtt
can2mqtt connects a CAN-Bus to an MQTT broker and vice versa. You create pairs of CAN-Bus IDs and MQTT topics and choose a conversion mode between them. 

Here you can see can2mqtt in action:
[![can2mqtt demo](screenshot.png)](https://asciinema.org/a/542608?autoplay=1)

## Installation
The latest compiled binary is available in the [releases](https://github.com/c3re/can2mqtt/releases/latest).
The binaries are statically linked and have no further dependencies. On a Raspberry you can install via:
```
wget https://github.com/c3re/can2mqtt/releases/download/v2.3.2/can2mqtt-v2.3.2-linux-arm -O can2mqtt
chmod +x can2mqtt
./can2mqtt
```

## Usage
The commandline parameters are the following:
 ```
 ./can2mqtt -f <can2mqtt.csv>  -m <mqtt-connectstring> [-v] -c <can-interface>
 ```
 
`<can2mqtt.csv>` is the configurations where your can2mqtt pairs are set.

`<mqtt-connectstring>` is a URL that is used to connect to your MQTT-Broker. username and password can be supplied like this: `tcp://user:pass@host:port`.

`<can-interface>` refers to the name of a [socket-can](https://www.kernel.org/doc/html/next/networking/can.html) interface on your computer that you want to bridge with can2mqtt. In case you have no idea what a socket-can interface is or how to configure it, read [how do i get a socket-can interface?](#how-do-i-get-a-socket-can-interface).

An additional `-v` flag can be passed to get verbose debug output.

Here is a full, working example:
```
./can2mqtt -f /etc/can2mqtt.csv -c can0 -m tcp://127.0.0.1:1883
```
## can2mqtt.csv
The file can2mqtt.csv has three columns. In the first column you need to specify the CAN-ID as a decimal number. In the second column you have to specify the convert-mode. You can find a list of available convert-modes below. In the last column you have to specify the MQTT-Topic. Each CAN-ID and each MQTT-Topic is allowed to appear only once in the whole file.

Example:
```
112,none,huette/all/a03/door/sensors/opened
113,2uint322ascii,huette/all/000/ccu/sensors/time
115,uint322ascii,huette/serverraum/000/filebitch/sensors/ftp_diskusage_percent
116,uint322ascii,huette/all/000/router/sensors/rx_bytes_s
117,uint322ascii,huette/all/000/router/sensors/tx_bytes_s
118,uint322ascii,huette/clubraum/000/ds18b20/sensors/temperatur
119,uint322ascii,huette/all/000/airmonitor/sensors/temp
120,uint322ascii,huette/all/000/airmonitor/sensors/hum
121,uint322ascii,huette/all/000/airmonitor/sensors/airq
122,uint322ascii,huette/all/000/airmonitor/sensors/pm2_5
123,uint322ascii,huette/all/000/airmonitor/sensors/pm10
```

Explanation for the 1st Line: For example our Doorstatus is published on the CAN-Bus every second with the CAN-ID 112 (decimal). can2mqtt will take everything that is published there and will push it through to mqtt-topic huette/all/a03/door/sensors/opened.

## convert-modes
Here they are:

| convertmode           | description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
|-----------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `none`                | does not convert anything. It just takes a bunch of bytes and hands it over to the other side. If you want to send strings, this will be your choice. If you have a mqtt payload that is longer than eight bytes, only the first eight bytes will be send via CAN.                                                                                                                                                                                                                                                                                                                                                                                                                                                        |
| `bytecolor2colorcode` | Converts an bytearray of 3 bytes to hexadecimal colorcode                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `pixelbin2ascii`      | This mode was designed to address colorized pixels. MQTT-wise you can insert a string like "<0-255> #RRGGBB" which will be converted to 4 byte on the CAN-BUS the first byte will be the number of the LED 0-255 and bytes 1, 2, 3 are the color of red, green and blue.                                                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `16bool2ascii`        | Interprets two bytes can-wise and publishes them as 16 boolean values to mqtt                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| `{i}[u]int{b}2ascii`  | `i` is the amount of instances of numbers in the CAN-Frame/the MQTT-message. Valid instance amounts are currently 1,2,4 and 8. Although other combinations are possible. Might be added in the future, if the need arises. `b` is the size of each number in bits. Supported values are 8,16,32 and 64. You can use either one unsigned or signed integers. This are all possible combinations:  `int82ascii`, `2int82ascii`, `4int82ascii`, `8int82ascii`, `int162ascii`, `2int162ascii`, `4int162ascii`, `int322ascii`, `2int322ascii`, `int642ascii`, `uint82ascii`, `2uint82ascii`, `4uint82ascii`, `8uint82ascii`, `uint162ascii`, `2uint162ascii`, `4uint162ascii`, `uint322ascii`, `2uint322ascii`, `uint642ascii` |


## Unidirectional Mode
Normally can2mqtt works in bidirectional mode, that means all messages from the CAN-bus are send to mqtt and vice versa. If you wish you can run can2mqtt in a unidirectional mode to only send messages from can to mqtt or only mqtt to can. To do so you have to use the flag `-d` with one fo the following settings:

| dirMode | effect                                                                          |
|---------|---------------------------------------------------------------------------------|
| 0       | bidirectional mode, messages will be send from can to mqtt and vice versa       |
| 1       | unidirectional mode, messages will only be send from CAN-Bus to mqtt broker     |
| 2       | unidirectional mode, messages will only be send from mqtt broker to the CAN-Bus |
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
use your convert-Mode in the `can2mqtt.csv` config-file. Now, use the file `src/convertmode/mymode.go` as a template for your own convertmode. Copy that file to `src/convertmode/<yournewmode>.go`. Now change all occurrences of "MyMode" with your preferred Name (Lets say `YourNewMode` in this example). Note that it has to start with an upper-case letter, so that it is usable outside of this package. Next you have to write three functions (implement the `ConvertMode` interface):
1. A conversion method from CAN -> MQTT: `ToMqtt(input can.Frame) ([]byte, error)`
2. A conversion method from MQTT -> CAN: `ToCan(input []byte) (can.Frame, error)`
3. A `String() string` method that reports the name of that convertmode. This method is used in some log-messages

Your almost done, the last step is to "register" your new convertmode. To do so add the following line to [`src/main.go#L72`](./src/main.go#L72)
```go
convertModeFromString[convertmode.<YourNewMode>{}.String()] = convertmode.<YourNewMode>{}
```

Now you can use your new convertmode in your `can2mqtt.csv` config File. Use the string that your return in the `String()` function as the name of the convertmode. In the `mymode.go` code this is `"mymode"`.

Good luck & happy hacking ✌
