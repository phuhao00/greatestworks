#!/usr/bin/env bash

gobin=$(go env GOBIN)

if test -z $gobin ; then
    gopath=$(go env GOPATH)
    if test -z $gopath; then
      gopath="$HOME/go"
    fi
    gobin="$gopath/bin"
fi

export PATH="$PATH:$gobin"

protoc=${which protoc}

if test -z $protoc; then
  printf "protoc not found"
  exit 1
fi

gengo=${which protoc-gen-go}
if test -z $gengo; then
  printf "gengo not found "
  exit  1
fi

exec protoc --go_out=. --go_opt=paths=source_relative,protofile=./pb ${1+"$@"}
