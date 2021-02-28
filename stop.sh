#!/bin/sh

/usr/bin/docker kill --signal="SIGINT" tentis > /dev/null
/usr/bin/docker stop -t 10 tentis > /dev/null
/usr/bin/docker rm tentis > /dev/null