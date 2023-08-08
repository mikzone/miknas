#!/bin/bash
set -e
webrootdir="$(readlink -f webroot)"
# export MIKNAS_DATABASE_DEBUG=1
export MIKNAS_DATABASE_PATH=$webrootdir/config/miknas.sqlite
export MIKNAS_CONFIG_DIR=$webrootdir/config
export MIKNAS_WORKSPACE=$webrootdir/workspace
echo "Workspace:" $MIKNAS_WORKSPACE
echo "Datebase:" $MIKNAS_DATABASE_PATH
cd ../server/example
if ! [ -d client ]; then
  cd ../../client
  quasar build
  cd ../server/example
  cp -r ../../client/dist/spa ./client
fi
go run main.go
