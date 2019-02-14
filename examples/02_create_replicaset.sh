#! /bin/bash
echo "Example: Creating a replica set"
read -p "Press ENTER to start"
echo""

echo "Cleaning cluster"
kubectl delete fooreplicaset --all
kubectl delete foo --all
echo ""

echo "Creating replica set"
kubectl apply -f config/samples/foo_v1_fooreplicaset.yaml
echo ""

kubectl get allfoo
