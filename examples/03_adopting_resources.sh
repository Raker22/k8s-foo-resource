#! /bin/bash
echo "Example: Adopting a resource"
read -p "Press ENTER to start"
echo""

echo "Cleaning cluster"
kubectl delete fooreplicaset --all
kubectl delete foo --all
echo ""

echo "Creating foo"
kubectl apply -f config/samples/foo_v1_foo.yaml
kubectl get foos
echo ""

echo "Creating foo replica set"
kubectl apply -f config/samples/foo_v1_fooreplicaset.yaml
echo ""

kubectl get allfoo
