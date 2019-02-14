#! /bin/bash
echo "Example: Existing resources when controllers start"
read -p "Stop controllers (press ENTER to continue)"
echo ""

echo "Cleaning cluster"
kubectl delete fooreplicaset --all
kubectl delete foo --all
echo ""

echo "Creating Resources"
kubectl apply -f config/samples/foo_v1_fooreplicaset.yaml
echo ""

kubectl get allfoo
echo ""

read -p "Start controllers now (press ENTER to continue)"
echo ""

kubectl get allfoo