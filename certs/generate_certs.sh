#!/bin/bash

# Create mqtt certificates ca.crt, server.crt, server.key client.crt, client.key

openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/C=US/ST=CA/L=San Francisco/O=My Company/CN=ca"

openssl genrsa -out server.key 2048
openssl req -new -out server.csr -key server.key -subj "/C=US/ST=CA/L=San Francisco/O=My Company/CN=localhost"
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650

openssl genrsa -out client.key 2048
openssl req -new -out client.csr -key client.key -subj "/C=US/ST=CA/L=San Francisco/O=My Company/CN=client"
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 3650

scp ca.crt server.crt server.key client.crt client.key root@10.0.2.249:/user
