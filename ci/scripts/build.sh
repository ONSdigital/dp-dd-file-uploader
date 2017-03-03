#!/bin/bash -eux

export BINPATH=$(pwd)/bin
export GOPATH=$(pwd)/go

pushd $GOPATH/src/github.com/ONSdigital/dp-dd-file-uploader
  go build -o $BINPATH/dp-dd-file-uploader
popd
