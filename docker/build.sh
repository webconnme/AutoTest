#!/bin/sh

IMAGES="at-dev"
for IMAGE in $IMAGES
do
	if [ -d $IMAGE ]
	then
		VERSION=$(cat $IMAGE/version)
		docker build -t webconn/$IMAGE:$VERSION $IMAGE && \
		docker build -t webconn/$IMAGE:latest $IMAGE
	fi
done
