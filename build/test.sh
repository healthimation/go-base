#!/bin/bash
set -e

ServiceName="goal-test"
verboseMode=""
runOpt=""
clean="false"
kubeContext="docker-desktop"

#handle args
while getopts ":vr:" opt; do
  case $opt in
    v)
      verboseMode="-v"
      ;;
    r)
      runOpt="-run (?i)$OPTARG"
      ;;
    c)
      clean="true"
      ;;
  esac
done

export SERVICE_NAME=$ServiceName
docker build -t $ServiceName-migrate -f Dockerfile.migrate .

echo -e "\033[0;35mSwitching k8s Context\033[0m"
# target the local docker k8s cluster
kubectl config use-context ${kubeContext}

if [ $clean = "true" ]
then 
  echo -e "\033[0;35mDestroying service stack and wiping database\033[0m"
  # deploy the service stack
  kubectl delete -f config/k8s/service/test || true
  rm -r "/tmp/k8s-${ServiceName}-db" || true
fi

kubectl apply -f config/k8s/service/test

echo -e "\033[0;35mWaiting on db..\033[0m"
# wait for db to be ready
kubectl wait --for=condition=ready pod -l app=${ServiceName}-db
sleep 2

echo -e "\033[0;35mRunning migrations...\033[0m"
# run migrations and wait
kubectl delete -f config/k8s/migrate/test || true
kubectl apply -f config/k8s/migrate/test
sleep 2

done=0
count=0
while [ $done = 0 ]
do
    mostRecentPod=$(kubectl get pods -l app=${ServiceName}-migrate --sort-by=.status.startTime | tail -n 1 | awk '{print $1}')
    kubectl logs -f $mostRecentPod
    
    if [[ $(kubectl get job ${ServiceName}-migrate -o=jsonpath='{.status.succeeded}') = '1' ]]
    then
        done=1
    elif [ "$count" -gt 30 ]
    then
        echo -e "\033[0;31mMigration Failed: exceeded max tries\033[0m"
        exit 1
    else
        let "++count"
    fi

    sleep 2
done

# lookup db host
if [ ! -n "${DATABASE_HOST+1}" ]; then
  dbport=$(kubectl describe svc ${ServiceName}-db | grep 'NodePort:' | awk '{print $3}' | cut -d/ -f 1)
  echo "Defaulting DATABASE_HOST based on kubectl lookup: localhost:$dbport"
  export DATABASE_HOST="localhost:$dbport"
fi


# default service host
if [ ! -n "${SERVICE_HOST+1}" ]; then
  port=$(kubectl describe svc ${ServiceName} | grep 'NodePort:' | awk '{print $3}' | cut -d/ -f 1)
  echo "Defaulting SERVICE_HOST to: localhost:${port}"
  export SERVICE_HOST="localhost:${port}"
fi

# find all go packages
packages="$(find src -type f -name "*.go" -exec dirname {} \; | sort | uniq)"
echo $packages
lintRet=0
vetRet=0
testRet=0
#loop through packages and test
for p in $packages
  do
    # golint if it is installed
    if golint 2>/dev/null; then
      echo "linting package $p"
      golint $p/*.go
      lintRet=$lintRet+$?
    fi

    # vet
    echo "running go vet on $p"
    go vet $p/*.go
    vetRet=$vetRet+$?

    # test
    echo "Running tests for $p"

    # make a tmp cover file then copy it to the right location for SublimeGoCoverage
    cover=$p/cover.out
    tmpcover=$(mktemp /tmp/tmp.XXXXXX)

    go test $verboseMode -coverprofile $tmpcover $runOpt "./$p"
    testRet=$testRet+$?

    sed 's/.*\///' $tmpcover > $cover
  done

# fail if any of the tests / vet / lint failed
exit $(($lintRet+$vetRet+$testRet))
