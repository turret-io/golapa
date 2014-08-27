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
if [ ! -d src ]; then
    echo "No source directory to copy"
    exit
fi
rsync -av src $MOUNT_POINT --exclude .git
cp -R src $MOUNT_POINT
docker build -t turretio/go-runtime docker/go-runtime
