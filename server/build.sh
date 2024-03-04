#!/bin/bash

# Set the target architecture
export GOARCH=arm
export GOARM=7

PRINTER_IP="10.0.2.249"

# Set the target operating system (if needed)
# export GOOS=linux

# Set the path to the Go compiler for ARM
export CC=arm-linux-gnueabihf-gcc

# If no dist then create it
if [ ! -d "dist" ]; then
    mkdir dist
fi

# check if prod flag is set
if [ "$1" == "prod" ]; then
    echo "Building for production"
    go build -o dist/kobra-server -ldflags "-s -w"
elif [ "$1" == "clean" ]; then
    echo "Cleaning up"
    rm -rf dist
    go clean
elif
    [ "$1" == "install" ]
then
    echo "Installing"
    # scp
    scp -r dist root@$PRINTER_IP:/opt/kobra
else
    echo "Building for development"
    go build -o dist/kobra-server
fi
# Clean up the environment variables
unset GOARCH
unset GOARM
unset CC
