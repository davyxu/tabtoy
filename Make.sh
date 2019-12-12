#!/usr/bin/env bash

export GOPROXY=https://goproxy.io

go build -v -o ./tabtoy.exe github.com/davyxu/tabtoy