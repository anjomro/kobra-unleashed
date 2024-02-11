#!/bin/bash

#One line
details="-subj /C=US/ST=CA/L=SanFrancisco/O=GlobalSecurity/OU=IT/CN=example.com"

#!/bin/bash

openssl genrsa -out ca.key 2048 && openssl req -x509 -new -nodes -key ca.key -sha256 -days 2028 -out ca.pem $details

openssl genrsa -out verification_cert.key 2048 && openssl req -new -key verification_cert.key -out verification_cert.csr $details && openssl x509 -req -in verification_cert.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out verification_cert.pem -days 2048 -sha256

openssl genrsa -out client.key 2048 && openssl req -new -key client.key -out client.csr $details && openssl x509 -req -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out client.pem -days 500 -sha256

mv ca.pem ca.crt
mv client.pem client.crt

# Md5sum ca.crt client.crt client.key verification_cert.pem verification_cert.key

md5sum ca.crt client.crt client.key verification_cert.pem verification_cert.key

# Create scp command to copy the files to the server
echo "scp ca.crt client.crt client.key verification_cert.pem verification_cert.key root@10.0.2.249:/user"
