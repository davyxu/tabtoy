#!/usr/bin/env bash
../../../../../../bin/tabtoy \
--mode=v2 \
--csharp_out=./csharp/Example/Config.cs \
--binary_out=./csharp/Example/Config.bin \
--lua_out=./lua/Config.lua \
--luaenumintvalue=true \
--go_out=./golang/table/table_gen.go \
--json_out=./golang/Config.json \
--cpp_out=./cpp/cpp/Config.h \
--combinename=Config \
--lan=zh_cn \
Globals.xlsx \
Sample.xlsx