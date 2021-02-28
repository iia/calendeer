#!/bin/sh

./stop.sh

/usr/bin/docker \
run \
--rm \
-p 5000:5000 \
--name tentis \
-i iia86/tentis:arm32v7-latest
