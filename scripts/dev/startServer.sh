#!/bin/bash
set -e
webrootdir="$(readlink -f ../../webroot)"
# export MIKNAS_DATABASE_DEBUG=1
export MIKNAS_DATABASE_PATH=$webrootdir/miknas.sqlite
export MIKNAS_WORKSPACE=$webrootdir/workspace
echo "Workspace:" $MIKNAS_WORKSPACE
echo "Datebase:" $MIKNAS_DATABASE_PATH
cd ../../server/example
go run main.go
