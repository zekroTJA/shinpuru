#!/bin/bash

set -e

IMAGE_NAME="zekro/shinpuru"

BRANCH=$TRAVIS_BRANCH
TAG=$TRAVIS_TAG

DTAG=""

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

IMGAE="$IMAGE_NAME:$DTAG"

docker build . -t $IMAGE

docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD

docker push $IMAGE