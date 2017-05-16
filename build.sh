#!/bin/bash

INITIAL_BASE_PATH=`pwd`

echo "Building with environment:"
echo "  GOROOT: `go env GOROOT`"
echo "  GOPATH: `go env GOPATH`"
echo "    GOOS: `go env GOOS`"
echo "  GOARCH: `go env GOARCH`"

GOPATH=`go env GOPATH`
#GOOS=solaris
#GOARCH=amd64
#export GOOS GOARCH

TARGET_BIN_FOLDER=$GOPATH/bin/solaris

if [ ! -x ${TARGET_BIN_FOLDER} ] ; then mkdir -p ${TARGET_BIN_FOLDER} ; fi

# TODO: There are most likely regular "go ways" to do the following .. but it does the trick for now ..
#
    
#go get .. all the dependencies

BUILD_PRODUCT_NAME=zfs_exporter
BUILD_PRODUCT_PATH=./cmd/${BUILD_PRODUCT_NAME}/${BUILD_PRODUCT_NAME}

cd ./cmd/${BUILD_PRODUCT_NAME} 
go build -v

cd ${INITIAL_BASE_PATH}


if [ -r ${BUILD_PRODUCT_PATH} ] ; then
    echo "Installing executable to: ${TARGET_BIN_FOLDER}/${BUILD_PRODUCT_NAME}"
    
    cp ${BUILD_PRODUCT_PATH} ${TARGET_BIN_FOLDER}/
fi

# END