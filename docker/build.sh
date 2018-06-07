#!/bin/bash -e

source conf.sh

function build {
    SERVICE=$1
    IMAGE_VERSION=$2
    PARENT_IMAGE=$3

    echo "########"
    echo "Building container $SERVICE:$IMAGE_VERSION PARENT:$PARENT_IMAGE"
    echo "########"

    IMAGE_NAME="coscale/$SERVICE:$IMAGE_VERSION"
    (docker images --format="{{.Repository}}:{{.Tag}}" | grep $IMAGE_NAME) && return

    docker pull $PARENT_IMAGE

    docker build --no-cache -t $IMAGE_NAME .
}

# Download the cli
./get-files.sh

# Build the cli image
build "cli" $VERSION "alpine:3.4"

# Remove files
./clean-files.sh
