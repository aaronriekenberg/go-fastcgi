#!/bin/bash -x

cd ~/go-fastcgi

systemctl --user stop go-fastcgi.service

git pull -v

time go build -x
RESULT=$?
if [ $RESULT -ne 0 ]; then
  echo "go build failed RESULT = $RESULT"
  exit $RESULT
fi

systemctl --user restart go-fastcgi.service
