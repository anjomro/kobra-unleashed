#!/bin/sh

# Start the mqtt server
mosquitto -c /app/mosquitto.conf &

sleep 1

# Start the kobra server
/app/server