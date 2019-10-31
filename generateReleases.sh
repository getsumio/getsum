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

tree -h builds
