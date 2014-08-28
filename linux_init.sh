#!/bin/bash

###
# 
# Setup 
# 
###

set -e 

CWD=`pwd`
VOLUME_PATH=$CWD
GOPROJ_PATH=$CWD/src/github.com/turret-io/golapa
SRC_DIR=golapa

if [ ! -d $GOPROJ_PATH ]; then
    mkdir -p $GOPROJ_PATH
fi

rsync -av $SRC_DIR/* $GOPROJ_PATH --exclude .git
docker run -v $VOLUME_PATH:/data --name golapa-build-volume busybox true
# Create .shared_volume symbolic link so we match other init scripts
ln -s $CWD .shared_volume
docker build -t turretio/go-runtime docker/go-runtime
