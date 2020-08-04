#!/bin/bash

#### Config Vars ####
# update these to reflect your service
ServiceName="<serviceName>"
kubeContext="docker-desktop"
awsSecretName="aws-auth"


# determine service path for volume mounting
CurrentDir=`pwd`
ServicePath="${CurrentDir/$BasePath/}"


# echo -e "\033[0;35mBuilding Service\033[0m"
# run the build
# docker run --rm -it -v `pwd`:"/go/src/$ServicePath" -w /go/src/$ServicePath golang:1.14 ./build/build.sh || { echo 'build failed' ; exit 1; }

echo -e "\033[0;35mBuilding Containers\033[0m"
# build the docker container with the new binary
docker build -t $ServiceName .
if [ $? -gt 0 ]; then exit 1; fi

docker build -t $ServiceName-migrate -f Dockerfile.migrate .


echo -e "\033[0;35mSwitching k8s Context\033[0m"
# target the local docker k8s cluster
kubectl config use-context ${kubeContext}

echo -e "\033[0;35mSetting up AWS credentials\033[0m"
# make sure aws auth is configure
if [ -z $AWS_ACCESS_KEY_ID ]
then
    echo -e "\033[0;31mNo AWS keys found in environment\033[0m"
    exit 1
fi
kubectl create secret generic ${awsSecretName} \
    --from-literal=AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID} \
    --from-literal=AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY} \
    --from-literal=AWS_SESSION_TOKEN=${AWS_SESSION_TOKEN} \
    --from-literal=AWS_SECURITY_TOKEN=${AWS_SECURITY_TOKEN} \
    --dry-run=client -o yaml | kubectl apply -f -

echo -e "\033[0;35mDeploying service and db stack\033[0m"
# deploy the service
kubectl delete deployment ${ServiceName} || true
kubectl apply -k config/k8s/service/local

echo -e "\033[0;35mWaiting on db..\033[0m"
# wait for db to be ready
kubectl wait --for=condition=ready pod -l app=${ServiceName}-db

echo -e "\033[0;35mRunning migrations...\033[0m"
# run migrations and wait
kubectl delete -k config/k8s/migrate/local || true
kubectl apply -k config/k8s/migrate/local
sleep 2

#kubectl get pods -l app=authentication-migrate --sort-by=.metadata.creationTimestamp -o jsonpath="{.items[-1].metadata.name}"
complete=0
failCount=0
while [ $complete = 0 ]
do
    
    if [[ $(kubectl get job ${ServiceName}-migrate -o=jsonpath='{.status.failed}') -gt 1 ]]
    then
        if [ "$failCount" -gt 30 ]
        then
            echo -e "\033[0;31mMigration Failed: exceeded max tries\033[0m"
            kubectl logs job/${ServiceName}-migrate
            exit 1
        else 
            echo "Waiting..."
            let "failCount++" 
        fi
    elif [[ $(kubectl get job ${ServiceName}-migrate -o=jsonpath='{.status.succeeded}') = '1' ]]
    then
        echo "Migration Completed"
        complete=1
    elif [[ $(kubectl get job ${ServiceName}-migrate -o=jsonpath='{.status.active}') = '1' ]]
    then
        echo "Migration Running..."
    else
        echo "unexpected migration job state!"
        kubectl get job ${ServiceName}-migrate -o yaml
        exit 1
    fi
    sleep 1
done

