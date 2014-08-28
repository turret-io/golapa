#!/bin/bash

###
# 
# Setup svendowideit/samba share for compiling the Go binaries on OS X
# 
###

set -e 

CWD=`pwd`
VOLUME_PATH=$CWD/._smbfs_share
MOUNT_POINT=$CWD/.shared_volume
GOPROJ_PATH=$MOUNT_POINT/src/github.com/turret-io/golapa
SRC_DIR=golapa
IP=`boot2docker ip`

if [ ! -d $VOLUME_PATH ]; then
    mkdir $VOLUME_PATH 
fi

if [ ! -d $MOUNT_POINT ]; then
    mkdir $MOUNT_POINT
fi

docker run -v $VOLUME_PATH:/data --name golapa-build-volume busybox true
docker run --rm -v /usr/local/bin/docker:/docker -v /var/run/docker.sock:/docker.sock svendowideit/samba golapa-build-volume
mount -t smbfs //guest:@$IP/data $MOUNT_POINT
if [ ! -d $SRC_DIR ]; then
    echo "No source directory to copy"
    exit
fi

if [ ! -d $GOPROJ_PATH ]; then
    mkdir -p $GOPROJ_PATH
fi

rsync -av $SRC_DIR/* $GOPROJ_PATH --exclude .git
docker build -t turretio/go-runtime docker/go-runtime
