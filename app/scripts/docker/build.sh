#!/bin/bash
set -e
# changeto app dir
cd ../../

# start build client
cd client
pwd
set -x
pnpm build
cd ../
set +x
pwd
set -x
rm -rf server/builded_client
cp -r client/dist server/builded_client

cd ..
docker build -t miknas -f app/scripts/docker/build.dockerfile .
docker image ls miknas
