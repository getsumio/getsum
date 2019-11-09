#!/bin/sh
VERSION=$(cat VERSION)
for GOOS in darwin linux windows openbsd freebsd ; do
    for GOARCH in 386 amd64; do
      GO_BIN_PATH="builds/$GOOS/$GOARCH"
        mkdir -p $GO_BIN_PATH
        EXTENSION=""
        if [[ "$GOOS" == "windows" ]];then
          EXTENSION=".exe"
        fi
        time GOOS=$GOOS GOARCH=$GOARCH go build -a -v -o ./$GO_BIN_PATH/getsum$EXTENSION ./cmd/getsum/main.go 
        TAR_FILE=$(pwd)/$GO_BIN_PATH/getsum-$GOOS-$GOARCH-$VERSION.tar.gz
        tar -czvf $TAR_FILE $GO_BIN_PATH
    done
done

ADDONS_PATH="builds/addons"
mkdir -p $ADDONS_PATH

pushd addons/chrome
sed -i 's/\"version\".*/\"version\": \"'"$VERSION"'\",/g' manifest.json
find . -exec zip ../../$ADDONS_PATH/chrome-$VERSION.zip
popd
pushd addons/firefox
sed -i 's/\"version\".*/\"version\": \"'"$VERSION"'\",/g' manifest.json
find . -exec zip ../../$ADDONS_PATH/firefox-$VERSION.zip
popd

tree -h builds
