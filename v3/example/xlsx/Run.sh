#!/bin/bash

# 默认设置代理, 国内专用
export GOPROXY=https://goproxy.io

go build -v -o ./tabtoy github.com/davyxu/tabtoy

./tabtoy -mode=v3 \
-index=Index.xlsx \
-go_out=../golang/table_gen.go \
-json_out=../json/table_gen.json \
-lua_out=../lua/table_gen.lua \
-csharp_out=../csharp/TabtoyExample/table_gen.cs \
-binary_out=../binary/table_gen.bin \
-java_out=../java/src/main/java/main/Table.java \
-package=main

cp ../json/table_gen.json ../java/cfg

rm -f tabtoy