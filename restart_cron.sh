#!/bin/sh

pgrep go-fastcgi-server > /dev/null 2>&1
if [ $? -eq 1 ]; then
  cd ~/go-fastcgi-server
  ./restart.sh > /dev/null 2>&1
fi
