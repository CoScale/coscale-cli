#!/bin/bash
# Tip: to remove '\r' use sed -i 's/\r//g' build.sh
export GOPATH=`pwd`
mkdir -p bin
go build -o bin/coscale-cli coscale
