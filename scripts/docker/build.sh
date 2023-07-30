#!/bin/bash
set -e

cd ../../client
pwd
set -x
quasar build
set +x

echo

cd ../
pwd
set -x
rm -rf server/example/client
cp -r client/dist/spa server/example/client
docker build -t miknas -f scripts/docker/build.dockerfile .
docker image ls miknas