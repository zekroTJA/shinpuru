#!/bin/bash

IMAGE_NAME="zekro/shinpuru"

BRANCH=$(git rev-parse --abbrev-ref HEAD)
TAG=$(git describe --tags --abbrev=0)

DTAG=""

set -e

echo "BRANCH: $BRANCH"
echo "TAG:    $TAG"

if ! [ -z $TAG ]; then
    DTAG=$TAG
else
    case "$BRANCH" in
        "dev")
            DTAG="canary"
            ;;
        "master")
            DTAG="latest"
            ;;
        *)
            exit 0
            ;;
    esac
fi

IMAGE="$IMAGE_NAME:$DTAG"

docker build . -t $IMAGE
docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
docker push $IMAGE