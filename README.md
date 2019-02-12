# Kubernetes Custom Resouce

Creating a custom resource with [Kubebuilder](http://book.kubebuilder.io)

## Project Creation

Initialize the project in `$GOPATH/src/...`

```
$ kubebuilder init --domain raker22.com
```

Create APIs

```
$ kubebuilder create api --group foo --version v1 --kind Foo
$ kubebuilder create api --group foo --version v1 --kind FooReplicaSet
```

## Configuration

* [Setup types](pkg/apis/foo/v1/README.md) with desired fields.
* [Update controllers](pkg/controller/README.md) with desired behavior.
* [Create sample resources](config/samples) to see how the controllers work.

## Running the Manager

* Connect to the kubernetes cluster you want to run the manager on.
* Run `make run` to start the manager on the cluster.
* In another terminal run `kubectl apply -f config/samples/${group}-${version}-${kind}.yaml` to create a sample resource.
  * `group`, `version`, and `kind` are the lower case group, version, and kind from `kubebulder create api` above.
