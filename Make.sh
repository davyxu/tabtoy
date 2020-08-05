#!/usr/bin/env bash

#export GOPROXY=https://goproxy.io

if [ "$1" == "" ] 
then
	echo "Usage: Make.sh version(like 1.0.0)"
	exit 1
fi

Version=$1

go clean -cache

export GOARCH=amd64
export GOOS=windows

go build -v -o ./tabtoy.exe github.com/davyxu/tabtoy

tar zcvf tabtoy-${Version}-win64.tar.gz ./tabtoy.exe

export GOARCH=amd64
export GOOS=linux

go build -v -o ./tabtoy github.com/davyxu/tabtoy
	
tar zcvf tabtoy-${Version}-linux-x86_64.tar.gz ./tabtoy

export GOARCH=amd64
export GOOS=darwin

go build -v -o ./tabtoy github.com/davyxu/tabtoy
	
tar zcvf tabtoy-${Version}-osx-x86_64.tar.gz ./tabtoy