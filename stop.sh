#!/bin/bash
serverPID=`ps -ef | grep job-finder | grep -v grep | awk -F' ' '{print $2}'`
kill -9 $serverPID
dbPID=`docker ps | grep postgre | awk -F' ' '{print $1}'`
docker kill $dbPID
