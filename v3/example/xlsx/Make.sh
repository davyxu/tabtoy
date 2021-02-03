#!/bin/bash

go build -v -o ./tabtoy github.com/davyxu/tabtoy

./tabtoy -mode=v3 \
-index=Index.xlsx \
-go_out=../golang/table_gen.go \
-json_out=../json/table_gen.json \
-json_dir=../jsondir \
-lua_out=../lua/table_gen.lua \
-lua_dir=../luadir \
-binary_dir=../binary \
-csharp_out=../csharp/TabtoyExample/table_gen.cs \
-binary_out=../binary/table_gen.bin \
-java_out=../java/src/main/java/main/Table.java \
-proto_out=../protobuf/table.proto \
-pbbin_out=../protobuf/all.pbb \
-pbbin_dir=../protobuf \
-package=main

if [ $? -ne 0 ] ; then
	read -rsp $'Errors occurred...\n' ; 
	exit 1 
fi

cp ../json/table_gen.json ../java/cfg

rm -f tabtoy