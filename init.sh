#!/bin/bash

if [ -z ${1+x} ]; then echo "service name is not set.\n sh init.sh <serviceName>\n"; exit 1; else echo "using service name '$1'"; fi

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"

echo "replacing <serviceName> with $1..."
find $DIR/build -type f -name '*.sh' -print0 | xargs -0 sed -i '' "s/<serviceName>/$1/g"
find $DIR/src/main -type f -name '*.go' -print0 | xargs -0  sed -i '' "s/<serviceName>/$1/g" 
find $DIR/src/server/serviceName -type f -name '*.go' -print0 | xargs -0  sed -i '' "s/<serviceName>/$1/g" 
find $DIR/src/serviceName -type f -name '*.go' -print0 | xargs -0  sed -i '' "s/<serviceName>/$1/g" 
find $DIR//config -type f -name '*.yaml' -print0 | xargs -0  sed -i '' "s/<serviceName>/$1/g" 
find $DIR/.github -type f -name '*.yml' -print0 | xargs -0  sed -i '' "s/<serviceName>/$1/g" 

echo "moving serviceName directories..."
mv $DIR/src/server/serviceName $DIR/src/server/$1
mv $DIR/src/serviceName $DIR/src/$1

echo "\nFinished.  Don't forget to remove this script update the README.md and run go mod.\n"
