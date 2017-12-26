#!/bin/bash
eval $(minikube docker-env)
GOOS=linux go build -i .
docker build . -t roundrobin:$1
kubectl patch deployment roundrobin -p '{"spec":{"template":{"spec":{"containers":[{"name":"roundrobin","image":"roundrobin:'$1'"}]}}}}'