#!/bin/bash
set -e

wrapdocker &  
sleep 5

docker login -e circleci@lavaboom.io -u $DOCKER_USER -p $DOCKER_PASS https://registry.lavaboom.io
docker build -t registry.lavaboom.io/lavaboom/$CONTAINER_NAME .
docker push registry.lavaboom.io/lavaboom/$CONTAINER_NAME

start-stop-daemon --stop --pidfile "/var/run/docker.pid"