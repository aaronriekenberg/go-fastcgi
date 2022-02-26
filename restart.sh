#!/bin/sh

KILL_CMD=pkill

$KILL_CMD go-fastcgi-server

sleep 2

export PATH=${HOME}/bin:$PATH

nohup ./go-fastcgi-server 2>&1 > output &
