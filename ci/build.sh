#!/bin/bash

# STATICS
BUILDPATH="./bin"
BUILDNAME="shinpuru"
#########

SQLLDFLAGS=$(bash ./scripts/getsqlschemes.bash)

TAG=$(git describe --tags)
if [ "$TAG" == "" ]; then
    TAG="untagged"
fi

COMMIT=$(git rev-parse HEAD)

if [ ! -d $BUILDPATH ]; then
    mkdir $BUILDPATH
fi

BUILDS=( \
    'linux;arm' \
    'linux;amd64' \
    'windows;amd64' \
    'darwin;amd64' \
)

curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
dep ensure

for BUILD in ${BUILDS[*]}; do

    IFS=';' read -ra SPLIT <<< "$BUILD"
    OS=${SPLIT[0]}
    ARCH=${SPLIT[1]}
    BINARY=${BUILDPATH}/${BUILDNAME}_${OS}_$ARCH

    echo "Building ${OS}_$ARCH..."
    (env GOOS=$OS GOARCH=$ARCH \
        go build \
            -o ${BUILDPATH}/${BUILDNAME}_${OS}_$ARCH \
            -ldflags " \
                -X github.com/zekroTJA/shinpuru/internal/util.AppVersion=$TAG \
                -X github.com/zekroTJA/shinpuru/internal/util.AppCommit=$COMMIT \
                -X github.com/zekroTJA/shinpuru/internal/util.Release=TRUE \
                $SQLLDFLAGS" \
                ./cmd/shinpuru)
            

    if [ "$OS" = "windows" ]; then
        mv ${BUILDPATH}/${BUILDNAME}_windows_$ARCH $BUILDPATH/${BUILDNAME}_windows_${ARCH}.exe
        BINARY=$BUILDPATH/${BUILDNAME}_windows_${ARCH}.exe
    fi
done

sha256sum $BUILDPATH/* | tee $BUILDPATH/sha256sums

echo "Exporting commands manual..."
go run ./cmd/cmdman -o ./docs/commandsManual.md

wait