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
# Create .shared_volume symbolic link so we match other init scripts
ln -s $CWD .shared_volume
docker build -t turretio/go-runtime docker/go-runtime
