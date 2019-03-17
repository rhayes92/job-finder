#!/bin/bash 
./stop.sh
docker run  -p 5432:5432 -e POSTGRES_PASSWORD=admin -d postgres
sleep 10
pushd ~/goproj/src/job-finder
go install
job-finder
