#!/bin/bash

# This script will be executed by CI service over SSH
# on your remote server where shinpuru Canary is
# running on.
# It will clone the dev branch of the repository,
# build the binary, push it to the binaries location
# and restarts the shinpuru deamon.
# So your Canary server will always be on the latest
# dev build when you push into dev branch.

#------------------- SETTINGS ----------------------
# Default branch which will be cloned if no
# command line argument was passed
DEF_BRANCH=dev
# Destination where to move the resutlt binary to
BIN_DEST=~/servers/shinpuruCanary/shinpuruCanary
# The command to restart your shinpuru deamon
RESTART_CMD="pm2 restart shinpuruCanary"
#---------------------------------------------------

export GOPATH=~/TMPGOPATH

BRANCH=$1
[ -z $1 ] && BRANCH=$DEF_BRANCH

echo "Creating gopath $GOPATH..."
mkdir -p $GOPATH && {
    cd $GOPATH

    echo "Cloning repository..."
    git clone -b $BRANCH --depth 1 \
        https://github.com/zekroTJA/shinpuru.git \
        src/github.com/zekroTJA/shinpuru           && \
    cd src/github.com/zekroTJA/shinpuru            && \
    echo "Building binary..."                      && \
    bash scripts/build.sh                          && \
    echo "Moving binary to $BIN_TEST..."           && \
    mv -f shinpuru $BIN_DEST                       && \
    $RESTART_CMD &> /dev/null
    echo "Cleaning up..."
    rm -r -f $GOPATH
} || {
    echo "failed creating build directory. Aborting..."
}