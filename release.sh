#!/usr/bin/env bash

BASE_PATH=/usr/local/src/kube/clap
TARGET_PATH=/usr/kube/clap

rm -rf ${BASE_PATH}/kube/server/bin
mkdir -p ${BASE_PATH}/kube/server/bin
docker run -it --rm \
      -e GO111MODULE=on \
      -v ${BASE_PATH}:${TARGET_PATH} \
      golang:1.18-bullseye sh -c "cd ${TARGET_PATH} && go build -o kube/server/bin/clap main.go"
chmod +x ${BASE_PATH}/kube/server/bin/clap

cd ${BASE_PATH}/web
rm -rf ${BASE_PATH}/kube/server/ui
rm -rf build
yarn build
mv ${BASE_PATH}/web/build ${BASE_PATH}/kube/server/ui

cd ${BASE_PATH}/kube/server
bash build.sh
