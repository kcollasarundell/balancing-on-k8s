#!/bin/bash
eval $(minikube docker-env)
GOOS=linux go build -i .
docker build . -t consistenthash:$1
kubectl patch deployment consistenthash -p '{"spec":{"template":{"spec":{"containers":[{"name":"consistenthash","image":"consistenthash:'$1'"}]}}}}'