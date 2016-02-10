#!/bin/bash
export GOPATH=`pwd`
mkdir -p bin
CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w' -o bin/coscale-cli coscale
