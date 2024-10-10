#!/bin/bash

set -e

SCRIPT_DIR="$( dirname -- "$BASH_SOURCE"; )";

cd ${SCRIPT_DIR}/../../

tinygo build -target wasm -tags "purego noasm" -o node/test/testapp.wasm ./example/testapp
