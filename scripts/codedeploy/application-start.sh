#!/bin/bash

AWS_REGION=
CONFIG_BUCKET=
ECR_REPOSITORY_URI=
GIT_COMMIT=

INSTANCE=$(curl -s http://instance-data/latest/meta-data/instance-id)
CONFIG=$(aws --region $AWS_REGION ec2 describe-tags --filters "Name=resource-id,Values=$INSTANCE" "Name=key,Values=Configuration" --output text | awk '{print $5}')

(aws s3 cp s3://$CONFIG_BUCKET/dp-dd-file-uploader/$CONFIG.asc . && gpg --decrypt $CONFIG.asc > $CONFIG) || exit $?

source $CONFIG && docker run -d        \
  --env=AWS_REGION=$AWS_REGION         \
  --env=BIND_ADDR=$BIND_ADDR           \
  --env=KAFKA_ADDR=$KAFKA_ADDR         \
  --env=S3_URL=$S3_URL                 \
  --env=TOPIC_NAME=$KAFKA_TOPIC        \
  --env=UPLOAD_TIMEOUT=$UPLOAD_TIMEOUT \
  --name=dp-dd-file-uploader           \
  --net=$DOCKER_NETWORK                \
  --restart=always                     \
  $ECR_REPOSITORY_URI/dp-dd-file-uploader:$GIT_COMMIT
