# Controllers

## Default Controller

The default controller created by `kubebuilder create api` watches all resources of the given kind as well as deployments owned by the resources.
The reconciler for the resource creates a single deployment owned by the resource and updates the deployment's `Spec` to match a hardcoded template.

```go
// default controller add func
func add(mgr manager.Manager, r reconcile.Reconciler) error {
    // create controller for Kind
    c, err := controller.New("kind-controller", mgr, controller.Options{Reconciler: r})
    if err != nil {
        return err
    }
    
    // watch for changes to Kind
    err = c.Watch(&source.Kind{Type: &Kind{}}, &handler.EnqueueRequestForObject{})
    if err != nil {
        return err
    }
    
    // watch OtherKinds owned by Kind
    err = c.Watch(&source.Kind{Type: &OtherKind{}}, &handler.EnqueueRequestForOwner{
      IsController: true,
      OwnerType:    &Kind{},
    }
    if err != nil {
        return err
    }
    
    return nil
}
```

## Handlers

Handlers receive a resource from the system when it is changed and determine if the controller should respond to the change.
There are several predefined handlers provided by `sigs.k8s.io/controller-runtime/pkg/handler`.
The following are used by the default controller.

* `EnqueueRequestForObject` - Enqueues the target resource
* `EnqueueRequestForOwner` - Enqueues the owner of the target resource

Meaning that

* `kind-controller` responds to all `Kind` resource changes by reconciling the resource
* `kind-controller` responds to `OtherKind` resource changes by reconciling it's owner if the owner is a `Kind` resource

These two handlers work for most simple cases but [custom handlers](../handler) can be implemented.

## Reconciler
A `Reconciler` takes a resource that has changed and updates the system to match the desired state.
Update the `Reconcile` function for the reconciler `kubebuilder create api` creates to have the desired functionality.
