#!/usr/bin/env bash

set -e

GOOS=js GOARCH=wasm go build -o main.wasm main.go
python3 -m http.server
