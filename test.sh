#!/bin/bash
# Tip: to remove '\r' use sed -i 's/\r//g' test.sh

export GOPATH=`pwd`

go test -timeout 20m -a -v coscale/...
