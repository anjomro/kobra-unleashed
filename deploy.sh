#!/bin/bash

# Script to build and deploy the application to the printer

# Usage: ./deploy.sh user@<printer_ip>

# Check if the printer IP is provided
if [ -z "$1" ]; then
  echo "Please provide the ssh user and printer IP as an argument (user@<printer_ip>)"
  exit 1
fi

# Build the application
echo "Building the application"
./build.sh

if [ $? -ne 0 ]; then
  echo "Error building the application"
  exit 1
fi

# Deploy the application
echo "Deploying the application to the printer"

ssh "$1" "mkdir -p /opt/kobra"
# ./dist/backend -> /opt/kobra/server
scp dist/backend "$1":/opt/kobra/server

# Remove the existing frontend
ssh "$1" "rm -rf /www"
# ./dist/frontend -> /www/frontend
scp -r dist/frontend "$1":/www/

# ./server/.env -> /opt/kobra/.env
scp server/.env "$1":/opt/kobra/.env