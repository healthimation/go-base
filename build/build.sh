#!/bin/bash
build=$(pwd)
# Find the import path for this repo
name=$(grep rootPath vendor/vendor.json |head -n1| cut -d' ' -f2 |tr -d '",')
buildAll=false

#handle args
while getopts ":a" opt; do
  case $opt in
    a)
      buildAll=true
      echo "building all packages"
      ;;
  esac
done

# Build the service
if ! CGO_ENABLED=0 go build -a -ldflags "-s" -installsuffix cgo -o bin/app src/main/*.go
then
  exit 1
fi

# Also build the migrate binary (it must already be vendored in via Godeps!)
if $buildAll
then
  CGO_ENABLED=0 go build -a -ldflags "-s" -installsuffix cgo -o bin/migrate ${name}/vendor/github.com/fedyakin/migrate
fi
