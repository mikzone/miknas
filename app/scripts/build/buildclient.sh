#!/bin/bash
set -e

cd ../../client
pwd
set -x
pnpm build
set +x

echo

cd ../
pwd
set -x
rm -rf server/builded_client
mv client/dist server/builded_client
