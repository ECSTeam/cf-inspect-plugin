#!/bin/bash

GOOS=darwin go build -o inspect-plugin-osx
#GOOS=linux go build -o inspect-plugin-linux
#GOOS=windows GOARCH=amd64 go build -o inspect-plugin.exe
if [ $? != 0 ]; then
   printf "Error when executing compile\n"
   exit 1
fi
cf uninstall-plugin inspect
cf install-plugin -f ./inspect-plugin-osx
