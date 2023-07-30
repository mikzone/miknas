#!/bin/bash
set -e
cd ../../server/example
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64
echo "start building..."
go build -o ../../webroot/cgo_miknas_server -ldflags '-s -w' main.go

# changeto green color
echo -e "\033[32m"
echo "Static Build Succ!!!"
ls -lah $(readlink -f ../../webroot/cgo_miknas_server)
echo -e "\033[0m"
