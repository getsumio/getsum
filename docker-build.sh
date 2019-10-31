#!/bin/sh
VERSION=$(cat VERSION)
docker build -t getsum:$VERSION .
