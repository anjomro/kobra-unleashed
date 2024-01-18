#!/bin/bash


openssl genrsa -out ca.key 2048 && openssl req -x509 -new -nodes -key ca.key -sha256 -days 2028 -out ca.pem


openssl genrsa -out verification_cert.key 2048 && openssl req -new -key verification_cert.key -out verification_cert.csr && openssl x509 -req -in verification_cert.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out verification_cert.pem -days 2048 -sha256

openssl genrsa -out client.key 2048 && openssl req -new -key client.key -out client.csr && openssl x509 -req -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out client.pem -days 500 -sha256