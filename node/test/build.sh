#!/bin/bash

set -e

export PATH=/home/joseph/sdk/go1.25.0/bin/:$PATH

SCRIPT_DIR="$( dirname -- "$BASH_SOURCE"; )";

cd ${SCRIPT_DIR}/../../

which tinygo && \
  tinygo build -target wasm -tags "purego noasm" -o node/test/testapp-tinygo.wasm ./example/testapp

GOOS=js GOARCH=wasm go build -o node/test/testapp-go.wasm ./example/testapp
