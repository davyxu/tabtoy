#!/usr/bin/env bash
CurrDir=`pwd`
cd ../../../../../..
export GOPATH=`pwd`
go build -v -o ${GOPATH}/bin/tabtoy github.com/davyxu/tabtoy
cd ${CurrDir}

source Export.sh