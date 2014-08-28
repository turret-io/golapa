#!/bin/bash

#####
#
#  Build the Go binaries, docker container, and add the support files
#
#####

set -e

function usage {
    echo "Usage: ./build.sh [IMAGE NAME]"
}

[ -z "$1" ] && usage && exit 1 

if [ "$1" == "-h" ]; then
    usage
    exit 0
fi

docker run --rm --volumes-from=golapa-build-volume -e GOPATH=/data -w=/data turretio/go-runtime go get ./... 
docker run --rm --volumes-from=golapa-build-volume -e GOPATH=/data turretio/go-runtime go install github.com/turret-io/golapa/serve
docker run --rm --volumes-from=golapa-build-volume -e GOPATH=/data turretio/go-runtime go install github.com/turret-io/golapa/worker

mkdir -p docker/golapa/tmp

# Copy compiled binaries into directory where we'll be building the docker image
cp -R .shared_volume/bin docker/golapa/tmp/bin
cp -R static docker/golapa/tmp/static 
cp -R templates docker/golapa/tmp/templates

# Build the image from the Dockerfile in docker/golapa/
docker build -t $1 docker/golapa

# Remove the tmp dir
rm -r docker/golapa/tmp
