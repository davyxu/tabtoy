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
-pragma=BuiltinPragma.xlsx \
-go_out=buildintypes_gen.go \
-json_out=BuiltinData.json

# json转go代码嵌入tabtoy
JSONDATAFILE=jsondata_gen.go
echo "package table" > ${JSONDATAFILE}
echo "const builtinJson = \`" >> ${JSONDATAFILE}
cat BuiltinData.json >> ${JSONDATAFILE}
echo "\`" >> ${JSONDATAFILE}
gofmt -s -w ${JSONDATAFILE}