#!/bin/sh

pgrep go-fastcgi > /dev/null 2>&1
if [ $? -eq 1 ]; then
  cd ~/go-fastcgi
  ./restart.sh > /dev/null 2>&1
fi
