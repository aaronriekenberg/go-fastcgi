#!/bin/sh

KILL_CMD=pkill

$KILL_CMD go-fastcgi

sleep 2

export PATH=${HOME}/bin:$PATH

CONFIG_FILE="./configfiles/$(uname | tr '[:upper:]' '[:lower:]')-config.json"

nohup ./go-fastcgi $CONFIG_FILE > output 2>&1 &
