#!/bin/bash

set -e

IMAGE_NAME="zekro/shinpuru"

BRANCH=$TRAVIS_BRANCH
TAG=$TRAVIS_TAG

echo "TRAVIS_TAG '$TRAVIS_TAG'"
echo "TRAVIS_BRANCH '$TRAVIS_BRANCH'"

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

IMAGE="$IMAGE_NAME:$DTAG"

echo "IMAGE '$IMAGE'"

docker build . -t $IMAGE

docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD

docker push $IMAGE