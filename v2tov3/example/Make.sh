#!/bin/bash

CURR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`
cd ${CURR}

go build -v -o ${GOPATH}/bin/tabtoy github.com/davyxu/tabtoy

InputTableDir=../../v2/example
OutputTableDir=.

${GOPATH}/bin/tabtoy -mode=v2tov3 \
-upout=${OutputTableDir} \
${InputTableDir}/Globals.xlsx \
${InputTableDir}/Sample.xlsx


${GOPATH}/bin/tabtoy -mode=v3 \
-index=Index.xlsx \
-go_out=golang_gen.go \
-json_out=json_gen.json \
-package=example