#!/bin/bash

source scripts/getsqlschemes.bash

TAG=$(git describe --tags)
if [ "$TAG" == "" ]; then
    TAG="untagged"
fi

COMMIT=$(git rev-parse HEAD)

echo "Getting dependencies..."
go get -v -t ./...

echo "Building..."
go build -ldflags " \
    -X github.com/zekroTJA/shinpuru/internal/util.AppVersion=$TAG \
    -X github.com/zekroTJA/shinpuru/internal/util.AppCommit=$COMMIT \
    -X github.com/zekroTJA/shinpuru/internal/util.Release=TRUE \
    $SQLLDFLAGS" \
    ./cmd/shinpuru

wait