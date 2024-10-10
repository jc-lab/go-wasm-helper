#!/bin/bash

set -e

SCRIPT_DIR="$( dirname -- "$BASH_SOURCE"; )";

cd ${SCRIPT_DIR}/../../

export GOROOT=$HOME/go/go1.23.1
export PATH=$GOROOT/bin:$PATH

tinygo build -target wasm -tags "purego noasm" -o node/test/testapp.wasm ./example/testapp
