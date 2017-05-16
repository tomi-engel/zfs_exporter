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
    
#go get golang.org/x/sys/unix
#go build -ldflags '-linkmode internal' -o $GOPATH/bin/solaris/smartos_exporter github.com/tomi-engel/smartos_exporter

BUILD_PRODUCT_NAME=zfs_exporter
BUILD_PRODUCT_PATH=./cmd/${BUILD_PRODUCT_NAME}/${BUILD_PRODUCT_NAME}

cd ./cmd/${BUILD_PRODUCT_NAME} 
go build -v

cd ${INITIAL_BASE_PATH}


# TODO: There is most likely a regular "go way" to do the following ..
#

if [ -r ${BUILD_PRODUCT_PATH} ] ; then
    echo "Installing executable to: ${TARGET_BIN_FOLDER}/${BUILD_PRODUCT_NAME}"
    
    cp ${BUILD_PRODUCT_PATH} ${TARGET_BIN_FOLDER}/
fi

# END