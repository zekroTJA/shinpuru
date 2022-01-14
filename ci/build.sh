#!/bin/bash

# STATICS
BUILDPATH="./bin"
BUILDNAME="shinpuru"
#########

TAG=$(git describe --tags)
if [ "$TAG" == "" ]; then
    TAG="untagged"
fi

COMMIT=$(git rev-parse HEAD)

if [ ! -d $BUILDPATH ]; then
    mkdir $BUILDPATH
fi

BUILDS=(
    'linux;arm'
    'linux;amd64'
    'windows;amd64'
    'darwin;amd64'
)

for BUILD in ${BUILDS[*]}; do

    IFS=';' read -ra SPLIT <<< "$BUILD"
    OS=${SPLIT[0]}
    ARCH=${SPLIT[1]}

    echo "Building ${OS}_$ARCH..."
    (env GOOS=$OS GOARCH=$ARCH \
        go build \
            -o ${BUILDPATH}/${BUILDNAME}_${OS}_$ARCH \
            -ldflags " \
                -X github.com/zekroTJA/shinpuru/internal/util.AppVersion=$TAG \
                -X github.com/zekroTJA/shinpuru/internal/util.AppCommit=$COMMIT \
                -X github.com/zekroTJA/shinpuru/internal/util.Release=TRUE" \
                ./cmd/shinpuru)


    if [ "$OS" = "windows" ]; then
        mv ${BUILDPATH}/${BUILDNAME}_windows_$ARCH $BUILDPATH/${BUILDNAME}_windows_${ARCH}.exe
    fi
done

cd web
npm ci
npx ng build --configuration production --output-path ../bin/web/dist/web
cd ..

mkdir deploy
tar -C ./bin -czvf ./deploy/build_assets.tar.gz .

cd deploy
sha256sum ./build_assets.tar.gz | tee sha256sum.txt
md5sum ./build_assets.tar.gz | tee md5sum.txt
cd ..

wait
