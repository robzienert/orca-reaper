#!/bin/bash -e
MAIN_PATH="cmd/orcareaper/main.go"
rm -rf build
mkdir build
env GOOS=linux GOARCH=amd64 go build -o build/orca-reaper-linux-amd64 $MAIN_PATH
env GOOS=darwin GOARCH=amd64 go build -o build/orca-reaper-darwin-amd64 $MAIN_PATH
