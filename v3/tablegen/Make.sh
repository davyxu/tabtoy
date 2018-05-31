#!/bin/bash

# 内建类型构成:
# BuiltinIndex做索引，FieldType,ErrorID做数据文件
# 生成表数据data_gen.json，转为go文件编译嵌入
# 生成表类型types_gen.go，编译嵌入
# BuiltinTypes表手动编辑转换为data.go，作为数据表在启动时使用

CURR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`
cd ${CURR}

go build -v -o ${GOPATH}/bin/tabtoy github.com/davyxu/tabtoy

# 导出内置的json文件
${GOPATH}/bin/tabtoy \
-mode=v3 \
-package=table \
-builtinsymbol \
-index=BuiltinIndex.xlsx \
-go_out=../table/types_gen.go \
-json_out=data_gen.json

# json转go代码嵌入tabtoy
JSONDATAFILE=../table/data_gen.go
echo "package table" > ${JSONDATAFILE}
echo "const coreConfig = \`" >> ${JSONDATAFILE}
cat data_gen.json >> ${JSONDATAFILE}
echo "\`" >> ${JSONDATAFILE}
gofmt -s -w ${JSONDATAFILE}