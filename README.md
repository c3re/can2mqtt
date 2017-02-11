# can2mqtt


can2mqtt is a small piece of software written in Go. Its purpose is to be a bridge between a CAN-Bus and a MQTT-Broker. Those are completely different worlds but they have similiaritys in the way they are built. I think i don't have to speak about the differences so i will just pick up the similiarities: In the CAN-world you have so called CAN-Frames. Each CAN-Frame can contain up to eight bytes of payload and CAN-Frame has an ID. In the MQTT-world you have topics and messages. Each message has a specific topic. As you can see it should be possible to map CAN-IDs to MQTT-Topics and their respective payload to messages. That's what this little programm does.

## Installation
You need to have go installed. After that you should be able to get can2mqtt to run with the following commands:
```
$ mkdir go
$ export GOPATH=go
$ go get github.com/c3re/can2mqtt
```
After that you should have a runnable binary under $GOPATH/bin/ called can2mqtt.
 
## Usage
The commandline parameters are the following:
 ```
 can2mqtt -f <can2mqtt.csv> -c <can-interface> -m <mqtt-connectstring> [-v]
 ```
 
Where can2mqtt.csv is the file for the configuration of can and mqtt pairs, can-interface is a socketcan interface and mqtt-connectstring is string that is accepted by the eclipse paho mqtt client. An additional -v flag can be passed to get verbose debug output. Here an example that runs on our Raspberry Pi @c3RE:
```
can2mqtt -f /etc/can2mqtt.csv -c can0 -m tcp://127.0.0.1:1883
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
### none
does not convert anything. It just takes a bunch of bytes and hands it over to the other side. If you want to send strings, this will be your choice.
### uint82ascii / uint162ascii / uint322ascii / uint642ascii 
On the can2mqtt way it takes 1, 2, 4 or 8 byte and interprets it as an uint of that size and parses it to a human readable string for the mqtt side. The other way round this convert motde takes an int in a string representation and sends out an array of bytes representing that number (little-endian)
### 2uint322ascii
This one is a bit special but all it does is that it takes 8 bytes from the CAN-Bus and parses two uint32s out of it and sends them in a string representation to MQTT. The two numbers are seperated with a simple space(" "). MQTT2CAN-wise it takes two string representations of numbers and converts them to 8 bytes representing them as 2 uint32.
