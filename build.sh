#!/usr/bin/env bash

cd web
rm -rf build
yarn build

cd ../
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

rm -rf kube/server/bin kube/server/ui
mkdir -p kube/server/bin
mv main kube/server/bin/clap
chmod +x kube/server/bin/clap
mv web/build kube/server/ui

cd kube/server
bash build.sh

exec "$@"