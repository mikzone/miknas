#!/bin/bash
set -e
webrootdir="$(readlink -f webroot)"
export MIKNAS_DATABASE_DEBUG=1
export MIKNAS_DATABASE_PATH=$webrootdir/config/miknas.sqlite
export MIKNAS_CONFIG_DIR=$webrootdir/config
export MIKNAS_WORKSPACE=$webrootdir/workspace
export MIKNAS_SECRET_KEY=goodgood
echo "Workspace:" $MIKNAS_WORKSPACE
echo "Datebase:" $MIKNAS_DATABASE_PATH
cd ../server
if ! [ -d client ]; then
  cd ../client
  pnpm build
  cd ../server
  mv ../client/dist ./client
fi
export CGO_ENABLED=0
go run main.go
