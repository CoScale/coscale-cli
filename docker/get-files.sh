#!/bin/bash -e

source conf.sh

if [ -e ./files ]; then
	rm -Rf ./files
fi
mkdir ./files

cd ./files
	wget https://github.com/CoScale/coscale-cli/releases/download/$VERSION/coscale-cli
cd -
