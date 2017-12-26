#!/bin/bash
eval $(minikube docker-env)
GOOS=linux go build -i .
docker build . -t clusterip:$1
kubectl patch deployment cluster-ip -p '{"spec":{"template":{"spec":{"containers":[{"name":"cluster-ip","image":"clusterip:'$1'"}]}}}}'