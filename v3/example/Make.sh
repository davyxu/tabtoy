#!/bin/bash

CURR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`
cd ${CURR}

go build -v -o ${GOPATH}/bin/tabtoy github.com/davyxu/tabtoy

${GOPATH}/bin/tabtoy -mode=v3 \
-builtinsymbol=../table/BuiltinTypes.xlsx \
-index=Index.xlsx \
-go_out=./golang/golang_gen.go \
-json_out=json_gen.json \
-package=main