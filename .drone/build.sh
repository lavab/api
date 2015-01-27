#!/bin/bash
set -e

ensure_loop(){
  num="$1"
  dev="/dev/loop$num"
  if test -b "$dev"; then
    echo "$dev is a usable loop device."
    return 0
  fi

  echo "Attempting to create $dev for docker ..."
  if ! mknod -m660 $dev b 7 $num; then
    echo "Failed to create $dev!" 1>&2
    return 3
  fi

  return 0
}

LOOP_A=$(losetup -f)
LOOP_A=${LOOP_A#/dev/loop}
LOOP_B=$(expr $LOOP_A + 1)

ensure_loop $LOOP_A
ensure_loop $LOOP_B

wrapdocker &  
sleep 5

docker login -e circleci@lavaboom.io -u $DOCKER_USER -p $DOCKER_PASS https://registry.lavaboom.io
docker build -t registry.lavaboom.io/lavaboom/$CONTAINER_NAME .
docker push registry.lavaboom.io/lavaboom/$CONTAINER_NAME

start-stop-daemon --stop --pidfile "/var/run/docker.pid"