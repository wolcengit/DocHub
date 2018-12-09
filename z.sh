#!/bin/bash

# last is 3
TAG="wolcen/dochub:2.1-Z$1"

docker build -t $TAG .
docker tag  $TAG  172.10.60.2/$TAG
docker push  172.10.60.2/$TAG

