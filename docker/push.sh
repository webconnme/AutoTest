#!/bin/sh

IMAGES="at-dev"
for IMAGE in $IMAGES
do
	if [ -d $IMAGE ]
	then
		VERSION=$(cat $IMAGE/version)
		docker push webconn/$IMAGE
	fi
done
