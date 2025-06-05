#!/bin/sh

minikube config set driver kvm2
minikube config set cpus 8
minikube config set memory 60000
minikube start --kubernetes-version=v1.32.5
