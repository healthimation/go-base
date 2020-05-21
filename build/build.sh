#!/bin/bash

# Build the service
if ! CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -mod=vendor -a -ldflags "-s" -installsuffix cgo -o bin/app src/main/*.go
then
  exit 1
fi
