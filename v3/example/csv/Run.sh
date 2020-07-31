#!/bin/bash

# 默认设置代理, 国内专用
export GOPROXY=https://goproxy.io

go build -v -o ./tabtoy github.com/davyxu/tabtoy

./tabtoy -mode=v3 \
-index=Index.csv \
-go_out=../golang/table_gen.go \
-json_out=../json/table_gen.json \
-jsontype_out=../jsontype/type_gen.json \
-lua_out=../lua/table_gen.lua \
-csharp_out=../csharp/TabtoyExample/table_gen.cs \
-binary_out=../binary/table_gen.bin \
-java_out=../java/src/main/java/main/Table.java \
-package=main


if [[ $? -ne 0 ]] ; then
	read -rsp $'Errors occurred...\n' ;
	exit 1
fi

cp ../json/table_gen.json ../java/cfg

rm -f tabtoy