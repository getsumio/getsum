#!/bin/sh
for GOOS in darwin linux windows openbsd freebsd ; do
    for GOARCH in 386 amd64; do
      GO_BIN_PATH="builds/$GOOS/$GOARCH"
        mkdir -p $GO_BIN_PATH
        EXTENSION=""
        if [[ "$GOOS" == "windows" ]];then
          EXTENSION=".exe"
        fi
        time GOOS=$GOOS GOARCH=$GOARCH go build -v -o ./$GO_BIN_PATH/getsum$EXTENSION ./cmd/getsum/main.go 
        tar -czvf ./$GO_BIN_PATH/getsum-$GOOS-$GOARCH.tar.gz ./$GO_BIN_PATH/
    done
done
