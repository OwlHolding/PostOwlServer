#!/bin/bash
if [ -z "$1" ]
  then
    echo "Pass the host address for the certificate"
    echo "For example: ./certgen.sh postowl.org"
else
    mkdir certs
    openssl req -newkey rsa:2048 -sha256 -nodes -keyout certs/key.pem -x509 -days 3650 -out certs/cert.pem -subj "/C=RU/ST=Moscow/L=Moscow/O=PostOwl/CN=$1"
fi