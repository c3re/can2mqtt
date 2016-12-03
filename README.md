# can2mqtt


can2mqtt is a small piece of software written in Go. Its purpose is to be a bridge between a CAN-Bus and a MQTT-Broker. Those are completely different worlds but they have similiaritys in the way they are built. I think i don't have to speak about the differences so i will just pick up the similiarities: In the CAN-world you have so called CAN-Frames. Each CAN-Frame can contain up to eight bytes of payload. Each CAN-Frame has an ID. In the MQTT-world you have topics and messages. Each message is published to a specific topic. As you can see it should be possible to copy frames/messages from one world to the other. That's what this little programm does.

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
 can2mqtt <can2mqtt.csv> <can-interface> <mqtt-connectstring> [-v]
 ```
 
Where can2mqtt.csv is the file for the configuration of can and mqtt pairs, can-interface is a socketcan interface and mqtt-connectstring is string that is accepted by the eclipse paho mqtt client. An additional -v flag can be passed to get verbose debug output.

## can2mqtt.csv
The file can2mqtt.csv has three columns. In the first column you need to specify the CAN-ID as a decimal number. In the second column you have to specify the convert-mode. You can find a list of available convert-modes below. In the last column you have to specify the MQTT-Topic. Each CAN-ID and each MQTT-Topic is allowed to appear only once in the whole file.

## convert-modes
Currently there is an epic amount of two convert-modes available

### bytes2ascii
#### CAN->MQTT
takes all the bytes and publishes them as string to mqtt
#### CAN<-MQTT
takes the string and publishes the first eight bytes to the CAN
### byte2dec
#### CAN->MQTT
takes the first byte of the mqtt frame interprets it as an decimal number (0-255) and publishes the number as a string to MQTT
#### CAN<-MQTT
takes a string and tries to parse out a number that fits one byte.
