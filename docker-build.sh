#!/bin/sh
VERSION=$(cat VERSION)
docker build --no-cache -t getsum:$VERSION .
