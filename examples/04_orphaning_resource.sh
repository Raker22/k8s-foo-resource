#! /bin/bash
echo "Example: Orphaning resources"
read -p "Press ENTER to start"
echo""

echo "Cleaning cluster"
kubectl delete fooreplicaset --all
kubectl delete foo --all
echo ""

echo "Creating resources"
kubectl apply -f config/samples/foo_v1_foo.yaml
kubectl apply -f config/samples/foo_v1_fooreplicaset.yaml
echo ""

kubectl get allfoo
echo ""

echo "Changing labels on foo"
kubectl get foo foo-sample -o yaml |\
  sed "s/app: sample/app: other/g" |\
  kubectl apply -f -
echo ""

kubectl get allfoo
echo ""

echo "Creating matching replica set"
sed "s/name: fooreplicaset-sample/name: fooreplicaset-other/g" config/samples/foo_v1_fooreplicaset.yaml |\
  sed "s/app: sample/app: other/g" |\
  sed "s/message: hello world/message: goodbye/g" |\
  sed "s/value: 42/value: 24/g" |\
  kubectl apply -f -
echo ""

kubectl get allfoo
