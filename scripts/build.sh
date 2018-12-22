#!/bin/bash

OS="linux"
ARCH="amd64"

(env GOOS=$OS GOARCH=$ARCH \
    go build -o shinpuru_${OS}_$ARCH)