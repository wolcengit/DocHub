#!/bin/bash

# last is 3
TAG="wolcen/dochub:2.1-Z$1"
ALIYUN_TAG="registry.cn-hangzhou.aliyuncs.com/$TAG"
sudo docker build -t $TAG .
sudo docker push $TAG
sudo docker tag $TAG $ALIYUN_TAG
sudo docker push $ALIYUN_TAG

