#! /bin/bash

echo "Deleting FooReplicaSets"
kubectl delete fooreplicasets --all

echo "Deleting Foos"
kubectl delete foos --all
