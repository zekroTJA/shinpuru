#!/bin/bash

TAG=$(git describe --tags)
if [ "$TAG" == "" ]; then
    TAG="untagged"
fi

COMMIT=$(git rev-parse HEAD)

echo "Getting dependencies..."
go get -v -t ./...

echo "Building..."
go build \
    -ldflags "-X main.ldAppVersion=$TAG -X main.ldAppCommit=$COMMIT"

wait