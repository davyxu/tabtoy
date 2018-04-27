#!/bin/bash

CURR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`
cd ${CURR}

go build -v -o ${GOPATH}/bin/tabtoy github.com/davyxu/tabtoy

${GOPATH}/bin/tabtoy -mode=v3 -symbol=ExampleType.xlsx -go_out=golang_gen.go -json_out=json_gen.json -package=example ExampleData.xlsx