#!/bin/bash

#One line
details="-subj /C=US/ST=CA/L=SanFrancisco/O=GlobalSecurity/OU=IT/CN=0.0.0.0"

openssl genrsa -out ca.key 2048
openssl req -new -x509 -days 3650 -key ca.key -out ca.crt $details

openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr $details
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 3650
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr $details
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.crt -days 3650

md5sum ca.crt server.crt server.key client.key client.crt

# Gen scp
echo "scp ca.crt server.crt server.key client.key client.crt root@10.0.2.249:/user"
