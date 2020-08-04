#!/bin/bash

serviceName="<serviceName>"

kubectl delete -k config/k8s/service/local
kubectl delete -k config/k8s/migrate/local
rm -r /tmp/k8s-${serviceName}-db

kubectl delete -f config/k8s/service/test
kubectl delete -f config/k8s/migrate/test
rm -r /tmp/k8s-${serviceName}-test-db
