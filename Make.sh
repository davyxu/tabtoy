#!/usr/bin/env bash

CURRDIR=`pwd`
cd ../../../..
export GOPATH=`pwd`
cd ${CURRDIR}

go build -v -o tabtoy.exe github.com/davyxu/tabtoy