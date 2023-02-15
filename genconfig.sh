#!/bin/bash
if [ -z "$1" ]
  then
    echo "Enter a name for the new configuration file"
    echo "For example: config.json"
else
    echo "package main; func main() { GenerateConfig(\"$1\") }" > postowlconfigtempgenerator.go
    make config
    rm -f postowlconfigtempgenerator.go
fi