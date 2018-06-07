#!/bin/bash

export VERSION=$(curl -s -X GET https://api.github.com/repos/CoScale/coscale-cli/releases/latest | grep -Po '(?<="tag_name": ")[^"]+(?=")')

