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

| convertmode           | description                                                                                                                                                                                                                                                            |
|-----------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `none`                | does not convert anything. It just takes a bunch of bytes and hands it over to the other side. If you want to send strings, this will be your choice. If you have a mqtt payload that is longer than eight bytes, only the first eight bytes will be send via CAN.     |
| `bytecolor2colorcode` | Converts an bytearray of 3 bytes to hexadecimal colorcode                                                                                                                                                                                                              |
| `pixelbin2ascii`      | This mode was designed to adress colorized pixels. MQTT-wise you can insert a string like "<0-255> #RRGGBB" wich will be converted to 4 byte on the CAN-BUS the first byte will be the number of the LED 0-255 and bytes 1, 2, 3 are the color of red, green and blue. |
| `16bool2ascii`        | Interprets two bytes can-wise and publishes them as 16 boolean values to mqtt                                                                                                                                                                                          |
| *uint*                |                                                                                                                                                                                                                                                                        |
| `uint82ascii`         | one uint8 in the CAN-Frame to one uint8 as string in the mqtt payload                                                                                                                                                                                                  |
| `4uint82ascii`        | four uint8 in the CAN-Frame to four uint8 as string seperated by spaces in the mqtt payload.                                                                                                                                                                           |
| `8uint82ascii`        | eight uint8 in the CAN-Frame to eight uint8 as string seperated by spaces in the mqtt payload.                                                                                                                                                                         |
| `uint162ascii`        | one uint16 in the CAN-Frame to one uint16 as string in the mqtt payload                                                                                                                                                                                                |
| `4uint162ascii`       | four uint16 in the CAN-Frame to four uint16 as string seperated by spaces in the mqtt payload.                                                                                                                                                                         |
| `uint322ascii`        | one uint32 in the CAN-Frame to one uint32 as string in the mqtt payload                                                                                                                                                                                                |
| `2uint322ascii`       | two uint32 in the CAN-Frame to two uint32 as string seperated by spaces in the mqtt payload.                                                                                                                                                                           |
| `uint642ascii`        | one uint64 in the CAN-Frame to one uint64 as string in the mqtt payload                                                                                                                                                                                                |
| *int*                 |                                                                                                                                                                                                                                                                        |
| `int82ascii`          | one int8 in the CAN-Frame to one int8 as string in the mqtt payload                                                                                                                                                                                                    |
| `4int82ascii`         | four int8 in the CAN-Frame to four int8 as string seperated by spaces in the mqtt payload.                                                                                                                                                                             |
| `8int82ascii`         | eight int8 in the CAN-Frame to eight int8 as string seperated by spaces in the mqtt payload.                                                                                                                                                                           |
| `int162ascii`         | one int16 in the CAN-Frame to one int16 as string in the mqtt payload                                                                                                                                                                                                  |
| `4int162ascii`        | four int16 in the CAN-Frame to four int16 as string seperated by spaces in the mqtt payload.                                                                                                                                                                           |
| `int322ascii`         | one int32 in the CAN-Frame to one int32 as string in the mqtt payload                                                                                                                                                                                                  |
| `2int322ascii`        | two int32 in the CAN-Frame to two int32 as string seperated by spaces in the mqtt payload.                                                                                                                                                                             |
| `int642ascii`         | one int64 in the CAN-Frame to one int64 as string in the mqtt payload                                                                                                                                                                                                  |
