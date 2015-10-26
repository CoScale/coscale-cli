#!/bin/bash
# Tip: to remove '\r' use sed -i 's/\r//g' build.sh
export GOPATH=`pwd`/cli
mkdir -p bin
CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o bin/coscale-cli coscale
