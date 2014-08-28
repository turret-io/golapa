#!/bin/bash

###
# 
# Setup 
# 
###

set -e 

CWD=`pwd`
VOLUME_PATH=$CWD

if [ ! -d $VOLUME_PATH ]; then
    mkdir $VOLUME_PATH 
fi

docker run -v $VOLUME_PATH:/data --name golapa-build-volume busybox true
docker build -t turretio/go-runtime docker/go-runtime
