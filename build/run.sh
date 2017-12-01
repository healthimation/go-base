#!/bin/bash

#### Config Vars ####
# update these to reflect your service
source `pwd`/build/vars.sh
BasePath="/home/core/dev/"
port="8080"
dbPort="5432"

# The config env file to pass to docker run
EnvFile="`pwd`/config/dev.env"
if [ ! -f $EnvFile ];
then
  echo "Cannot find ${EnvFile}"
  exit 1
fi

SecretEnvFile="`pwd`/config/secret.env"
IncludeSecret="--env-file `pwd`/config/secret.env"
if [ ! -f $SecretEnvFile ];
then
  IncludeSecret=""
fi

#handle args
while getopts ":p:" opt; do
  case $opt in
    p)
      port=$OPTARG
      echo "port overridden to $port"
      ;;
  esac
done

# determine service path for volume mounting
CurrentDir=`pwd`
ServicePath="${CurrentDir/$BasePath/}"

# run the build
docker run --rm -it -v `pwd`:"/go/src/$ServicePath" -w /go/src/$ServicePath healthymation/docker-godep:1.7.3 ./build/build.sh -a || { echo 'build failed' ; exit 1; }

# build the docker container with the new binary
docker build -t $ServiceName .
docker build -t $ServiceName-migrate -f Dockerfile.migrate .


# run the DB container
DBName="$ServiceName-db"
# only run the db if it isnt already running
running=$(docker ps -q -f "name=$DBName" -f "status=running" )
if [ "$running" == "" ]
  then
  docker rm $DBName &>/dev/null
  docker run --name $DBName -e POSTGRES_PASSWORD=password -e POSTGRES_DB=$ServiceName -P -d postgres
  sleep 10
fi

# Detect DB ip:port
HostIP=$COREOS_PUBLIC_IPV4
DbPublicPort=`docker ps | grep $DBName | cut -d : -f2 | cut -d - -f1`
# register db on the host ip:port for the service and test.sh
resp=$(curl -s -o /dev/null -w "%{http_code}" -X PUT -d '{"Datacenter": "vagrant", "Node": "dev-db", "Address": "'$HostIP'", "Service": {"Service": "'$DBName'", "Address": "'$HostIP'", "Tags": ["dev"], "Port": '$DbPublicPort'}}' http://$HostIP:8500/v1/catalog/register)
if [ "$resp" != "200" ]
then
  echo "Non-200 response adding DB to Consul ($resp)"
  exit 1
fi

# run the migrations
echo "Running migrations on postgres: $DBName:$DbPublicPort"
docker run --rm -ti -v $CurrentDir/build/migration:/migration $ServiceName-migrate /migrate -url "postgres://postgres:password@$DBName:$DbPublicPort/$ServiceName?sslmode=disable" -path /migration up


# Use the services location env provided by the host file if it exists
useServicesEnv="--env-file /etc/services.env"
if [ ! -f "/etc/services.env" ];
then
  echo "Warning! /etc/services.env not found, service to service communication may not work as expected."
  $useServicesEnv=""
fi

# run the container
docker run --rm -it -p $port:8080 -e PGHOST=$DBName -e PGPORT=$DbPublicPort --env-file $EnvFile $IncludeSecret $useServicesEnv -e SERVICE_NAME=$ServiceName $ServiceName
