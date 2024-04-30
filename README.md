# can2mqtt


can2mqtt is a small piece of software written in Go. Its purpose is to be a bridge between a CAN-Bus and a MQTT-Broker. Those are completely different worlds but they have similiaritys in the way they are built. I think i don't have to speak about the differences so i will just pick up the similiarities: In the CAN-world you have so called CAN-Frames. Each CAN-Frame can contain up to eight bytes of payload and CAN-Frame has an ID. In the MQTT-world you have topics and messages. Each message has a specific topic. As you can see it should be possible to map CAN-IDs to MQTT-Topics and their respective payload to messages. That's what this little programm does.

Here you can see can2mqtt in action:
[![can2mqtt demo](screenshot.png)](https://asciinema.org/a/542608?autoplay=1)

## Installation
can2mqtt is written in Go and static linked binaries are available [here](https://github.com/c3re/can2mqtt/releases/latest).
can2mqtt has no further dependencies. On a Raspberry for example it should be enough to run:
```
wget https://github.com/c3re/can2mqtt/releases/download/v1.3.0/can2mqtt-v1.3.0-linux-arm -O can2mqtt
chmod +x can2mqtt
./can2mqtt
```

## Usage
The commandline parameters are the following:
 ```
 ./can2mqtt -f <can2mqtt.csv> -c <can-interface> -m <mqtt-connectstring> [-v]
 ```
 
Where can2mqtt.csv is the file for the configuration of can and mqtt pairs, can-interface is a socketcan interface and mqtt-connectstring is string that is accepted by the eclipse paho mqtt client. An additional -v flag can be passed to get verbose debug output. Here an example that runs on our Raspberry Pi @c3RE:
```
./can2mqtt -f /etc/can2mqtt.csv -c can0 -m tcp://127.0.0.1:1883
```
## can2mqtt.csv
The file can2mqtt.csv has three columns. In the first column you need to specify the CAN-ID as a decimal number. In the second column you have to specify the convert-mode. You can find a list of available convert-modes below. In the last column you have to specify the MQTT-Topic. Each CAN-ID and each MQTT-Topic is allowed to appear only once in the whole file.

Here again the example from the Pi@c3RE:

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

Explanation for the 1st Line: For example our Doorstatus is published on the CAN-Bus every second with the CAN-ID 112 (decimal). can2mqtt will take everything thats published there and will push it through to mqtt-topic huette/all/a03/door/sensors/opened.

## convert-modes
Here they are:

| interfaced | convertmode           | description                                                                                                                                                                                                                                                                                                                                              |
|------------|-----------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| X          | `none`                | does not convert anything. It just takes a bunch of bytes and hands it over to the other side. If you want to send strings, this will be your choice. If you have a mqtt payload that is longer than eight bytes, only the first eight bytes will be send via CAN.                                                                                       |
| X          | `16bool2ascii`        | Interprets two bytes can-wise and publishes them as 16 boolean values to mqtt                                                                                                                                                                                                                                                                            |
| X          | `uint82ascii `        | On the can2mqtt way it takes 1 byte and interprets it as an uint8 and converts it to a string for the mqtt side. The other way round this convert mode takes an uint8 in a string representation and sends out one byte representing that number (little-endian)                                                                                         |
| X          | `uint162ascii `       | On the can2mqtt way it takes 2 bytes and interprets them as an uint16 and converts it to a string containing that number for the mqtt side. The other way round this convert mode takes an uint16 in a string representation and sends out two bytes representing that number (little-endian)                                                            |
| X          | `uint322ascii `       | On the can2mqtt way it takes 4 bytes and interprets them as an uint32 and converts it to a string containing that number for the mqtt side. The other way round this convert mode takes an uint32 in a string representation and sends out four bytes representing that number (little-endian)                                                           |
| X          | `uint642ascii `       | On the can2mqtt way it takes 8 bytes and interprets them as an uint64 and converts it to a string containing that number for the mqtt side. The other way round this convert mode takes an uint64 in a string representation and sends out four bytes representing that number (little-endian)                                                           |
| X          | `2uint322ascii`       | This one is a bit special but all it does is that it takes 8 bytes from the CAN-Bus and parses two uint32s out of it and sends them in a string representation to MQTT. The two numbers are seperated with a simple space(" "). MQTT2CAN-wise it takes two string representations of numbers and converts them to 8 bytes representing them as 2 uint32. |
|            | `4uint162ascii`       | Interprets eight bytes can-wise and publishes them as 4 uint16 seperated by a space to the mqtt side                                                                                                                                                                                                                                                     |
|            | `4int162ascii`        | Interprets eight bytes can-wise and publishes them as 4 int16 seperated by a space to the mqtt side                                                                                                                                                                                                                                                      |
|            | `4uint82ascii`        | Interprets four bytes (byte 0, 2, 4 and 6) can-wise and publishes them as 4 uint8 seperated by a space to the mqtt side                                                                                                                                                                                                                                  |
|            | `8uint82ascii`        | Interprets eight bytes (byte 0 to 7) can-wise and publishes them as eight uint8 seperated by a space to the mqtt side. The other way around it expects eight bytes seperated by a space and publishes them as eight bytes on the can-side.                                                                                                               |
|            | `bytecolor2colorcode` | Converts an bytearray of 3 bytes to hexadecimal colorcode                                                                                                                                                                                                                                                                                                |
|            | `pixelbin2ascii`      | This mode was designed to adress colorized pixels. MQTT-wise you can insert a string like "<0-255> #RRGGBB" wich will be converted to 4 byte on the CAN-BUS the first byte will be the number of the LED 0-255 and bytes 1, 2, 3 are the color of red, green and blue.                                                                                   |
