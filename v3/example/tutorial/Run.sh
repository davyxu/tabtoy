#!/usr/bin/env bash

if [[ ! -f "./tabtoy" ]]; then
     echo "请在https://github.com/davyxu/tabtoy/releases下载最新的tabtoy, 并放置于本目录"
     exit 1
fi

./tabtoy -mode=v3 -index=Index.xlsx -json_out=table_gen.json
