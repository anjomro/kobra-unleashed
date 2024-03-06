#!/bin/bash

# Build the project

# Built fontend

# Set current directory
ORIG_DIR=$(pwd)

# Before build check if folder dist exists and delete it
if [ -d "dist" ]; then
    rm -rf dist
fi

cd kobra-client

# Build the frontend vite

yarn build

# Check if successful build
if [ $? -ne 0 ]; then
    echo "Error building frontend"
    exit 1
fi

# Copy dist to original directory /dist/frontend

cp -r dist ../dist/frontend

# Build the backend

cd $ORIG_DIR

cd server

# Build the backend using build.sh

./build.sh

if [ $? -ne 0 ]; then
    echo "Error building backend"
    exit 1
fi

# Copy dist/kobra-server to original directory /dist/backend

cp -r dist/kobra-server ../dist/backend

# Return to original directory

cd $ORIG_DIR

echo "Build complete. Check /dist for the build files."
