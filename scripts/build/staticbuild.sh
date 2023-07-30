#!/bin/bash
set -e
cd ../../server/example
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
echo "start building..."
go build -o ../../webroot/static_miknas_server -ldflags '-s -w --extldflags "-static -fpic"' main.go

# changeto green color
echo -e "\033[32m"
echo "Static Build Succ!!!"
ls -lah $(readlink -f ../../webroot/static_miknas_server)
echo -e "\033[0m"
