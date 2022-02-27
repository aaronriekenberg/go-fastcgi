#!/bin/sh

KILL_CMD=pkill

$KILL_CMD go-fastcgi

sleep 2

export PATH=${HOME}/bin:$PATH

nohup ./go-fastcgi ./configfiles/config.json > output 2>&1 &
