#!/bin/bash -e

source conf.sh
docker build -t coscale/cli --build-arg VERSION=${VERSION} .
