#!/usr/bin/env bash

# Build
go build -C ../src -o ../e2e-test/can2mqtt || exit

# Setup Virtual CAN Interface
# TODO find something that does not require root privileges
sudo ip link add dev vcan0 type vcan
sudo ip link set vcan0 up

# Setup local MQTT Server
# TODO find something that does not require root privileges, and is more system independent, probably a docker container
sudo systemctl start mosquitto.service

# Clean & Create logging directory
rm -rf logs
mkdir logs

# Start can2mqtt
./can2mqtt -c vcan0 -f e2e-test.csv > logs/can2mqtt 2>&1 &

# Start can-logging
candump vcan0 > logs/can 2>&1 &

# Start mqtt logging
mosquitto_sub -h localhost -v -t 'e2e-test/#'  > logs/mqtt 2>&1 &

# Run tests
# Publish MQTT-Tests
mosquitto_pub -h localhost -m "test" -t e2e-test/none
mosquitto_pub -h localhost -m "0 0 1 0 1 0 1 1 0 0 1 1 1 1 0 1" -t e2e-test/16bool2ascii
mosquitto_pub -h localhost -m "75" -t e2e-test/uint82ascii
mosquitto_pub -h localhost -m "35000" -t e2e-test/uint162ascii
mosquitto_pub -h localhost -m "2000000001" -t e2e-test/uint322ascii
mosquitto_pub -h localhost -m "123470851232" -t e2e-test/uint642ascii
mosquitto_pub -h localhost -m "0 123441234" -t e2e-test/2uint322ascii
mosquitto_pub -h localhost -m "12 89 1234 4" -t e2e-test/4uint162ascii
mosquitto_pub -h localhost -m "-1234 42 1243 2" -t e2e-test/4int162ascii
mosquitto_pub -h localhost -m "12 2 21 2" -t e2e-test/4uint82ascii
mosquitto_pub -h localhost -m "1 2 3 4 5 6 7 8" -t e2e-test/8uint82ascii
mosquitto_pub -h localhost -m "#00ff00" -t e2e-test/bytecolor2colorcode
mosquitto_pub -h localhost -m "12 #00ff00" -t e2e-test/pixelbin2ascii

# Send CAN-Tests
cansend vcan0 "064#ABCD"
cansend vcan0 "065#ABCD"
cansend vcan0 "066#ABCD"
cansend vcan0 "067#ABCD"
cansend vcan0 "068#ABCD"
cansend vcan0 "069#ABCD"
cansend vcan0 "06a#ABCD"
cansend vcan0 "06b#ABCD"
cansend vcan0 "06c#ABCD"
cansend vcan0 "06d#ABCD"
cansend vcan0 "06e#ABCD"
cansend vcan0 "06f#ABCDEF"
cansend vcan0 "070#ABCD"

#
sleep 5 
# Check results
md5sum logs/can
md5sum logs/mqtt

# Cleanup
pkill can2mqtt
pkill candump
pkill mosquitto_sub
cat logs/*
rm -rf logs
rm can2mqtt
sudo ip link delete vcan0
sudo systemctl stop mosquitto.service
