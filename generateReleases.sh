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
        TAR_FILE=$(pwd)/$GO_BIN_PATH/getsum-$GOOS-$GOARCH-v$VERSION.tar.gz
        tar -czvf $TAR_FILE $GO_BIN_PATH
    done
done

ADDONS_PATH="builds/addons/chrome"
mkdir -p $ADDONS_PATH

pushd addons/chrome
sed -i 's/\"version\".*/\"version\": \"'"$VERSION"'\",/g' manifest.json
find -L . -exec zip ../../$ADDONS_PATH/chrome-v$VERSION.zip {} \;
popd
ADDONS_PATH="builds/addons/firefox"
mkdir -p $ADDONS_PATH
pushd addons/firefox
sed -i 's/\"version\".*/\"version\": \"'"$VERSION"'\",/g' manifest.json
find -L . -exec zip ../../$ADDONS_PATH/firefox-v$VERSION.zip {} \;
popd

tree -h builds


find builds/ -name '*tar.gz' -exec getsum -q {} \;
