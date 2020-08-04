#!/bin/bash

serviceName="goal"
kubeContext="docker-desktop"

echo -e "\033[0;35mSwitching k8s Context\033[0m"
# target the local docker k8s cluster
kubectl config use-context ${kubeContext}

kubectl delete -k config/k8s/service/local
kubectl delete -k config/k8s/migrate/local
rm -r /tmp/k8s-${serviceName}-db

kubectl delete -f config/k8s/service/test
kubectl delete -f config/k8s/migrate/test
rm -r /tmp/k8s-${serviceName}-test-db
