#!/bin/sh

####
#
#  Unmounts the build volume for OS X and stops the samba container 
#
####

umount .shared_volume
docker stop golapa-build-volume
docker rm golapa-build-volume
