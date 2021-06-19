#!/usr/bin/env bash
Version=3.1.2

export GOARCH=amd64
BuildSourcePackage="github.com/davyxu/tabtoy/build"

BuildBinary()
{
  set -e
  TargetDir=bin/"${1}"
  mkdir -p "${TargetDir}"
  export GOOS=${1}
  BuildTime=$(date -R)
  GitCommit=$(git rev-parse HEAD)
  VersionString="-X \"${BuildSourcePackage}.BuildTime=${BuildTime}\" -X \"${BuildSourcePackage}.Version=${Version}\" -X \"${BuildSourcePackage}.GitCommit=${GitCommit}\""

  go build -v -p 4 -o "${TargetDir}"/tabtoy -ldflags "${VersionString}" github.com/davyxu/tabtoy
  PackageDir=$(pwd)
  cd "${TargetDir}"
  tar zcvf "${PackageDir}"/tabtoy-${Version}-"${1}"-x86_64.tar.gz tabtoy
  cd "${PackageDir}"
}



BuildBinary windows
BuildBinary linux
BuildBinary darwin