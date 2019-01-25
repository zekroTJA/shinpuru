#!/bin/bash

# This script will be executed by CI service over SSH
# on your remote server where shinpuru is running on.
# It will clone the tag branch of the repository,
# which must be passed by the CI ($TRAVIS_TAG i.e.),
# build the binary, push it to the binaries location
# and restarts the shinpuru deamon.
# So your server will always be on the latest tag
# when you create a new one.

#------------------- SETTINGS ----------------------
# Destination where to move the resutlt binary to
BIN_DEST=~/servers/shinpuru/shinpuru_linux_amd64
# The command to restart your shinpuru deamon
RESTART_CMD="pm2 restart shinpuru"
#---------------------------------------------------

[ -z $1 ] && {
    echo "tag not set as argument"
    exit 1
}

export GOPATH=~/TMPGOPATH

echo "Creating gopath $GOPATH..."
mkdir -p $GOPATH && {
    cd $GOPATH

    echo "Cloning repository..."
    git clone -b $1 --depth 1 \
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