#!/bin/bash

CURR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`
cd ${CURR}

go build -v -o ${GOPATH}/bin/tabtoy github.com/davyxu/tabtoy

# 导出内置的json文件
${GOPATH}/bin/tabtoy \
-mode=v3 \
-package=table \
-builtinsymbol=BuiltinTypes.xlsx \
-index=BuiltinIndex.xlsx \
-go_out=types_gen.go \
-json_out=data_gen.json

# json转go代码嵌入tabtoy
JSONDATAFILE=data_gen.go
echo "package table" > ${JSONDATAFILE}
echo "const coreConfig = \`" >> ${JSONDATAFILE}
cat data_gen.json >> ${JSONDATAFILE}
echo "\`" >> ${JSONDATAFILE}
gofmt -s -w ${JSONDATAFILE}