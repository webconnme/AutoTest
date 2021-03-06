#!/bin/sh

IMAGES="at-dev"
for IMAGE in $IMAGES
do
	if [ -d $IMAGE ]
	then
		VERSION=$(cat $IMAGE/version)
		docker rmi webconn/$IMAGE:$VERSION 
		docker rmi webconn/$IMAGE:latest
	fi
done
