# Custom Handlers

The built in handlers work for most simple cases.
However, `EnqueueRequestForOwner` relies on a resource already having `OwnerReferences` to it's owners.
If a resource wants to control other resources that match a selector, the resource won't be notified when matching resources change unless they're already owned by the resource.
If the resource created the matching resources and set itself as their owner, everything will work as expected but if a resource is created by another controller or manually created, the owner resource won't be notified of changes.

For this project I created a slightly more sophisticated handler that is similar to `EnqueueRequestForOwner` but solves the owner reference problem.
If a resource doesn't have a controller resource the handler attempts to find one to adopt it.
The handler has a `GetOwner` function property that takes a resource and returns a resource that could be it's controller.
This functionality is closer to how built in resources such as deployments or replica sets behave.
