package handler

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var _ handler.EventHandler = &EnqueueRequestForController{}

var log = logf.KBLog.WithName("eventhandler").WithName("EnqueueRequestForController")

// EnqueueRequestForController enqueues Requests for the Owners of an object.  E.g. the object that created
// the object that was the source of the Event.
//
// If a ReplicaSet creates Pods, users may reconcile the ReplicaSet in response to Pod Events using:
//
// - a source.Kind Source with Type of Pod.
//
// - a handler.EnqueueRequestForController EventHandler with an OwnerType of ReplicaSet and IsController set to true.
type EnqueueRequestForController struct {
	// OwnerType is the type of the Owner object to look for in OwnerReferences.  Only Group and Kind are compared.
	OwnerType runtime.Object

	// GetOwners returns owners of the object
	// - object is the object to get owners for
	// - onlyOne if set means that only the first owner returned will be set as the objects owner
	GetOwner func(object metav1.Object) (metav1.Object, error)

	// Scheme is the
	Scheme *runtime.Scheme

	// groupKind is the cached Group and Kind from OwnerType
	groupKind schema.GroupKind
}

// Create implements EventHandler
func (e *EnqueueRequestForController) Create(evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if req, err := e.getOwnerReconcileRequest(evt.Meta); err == nil {
		q.Add(req)
	}
}

// Update implements EventHandler
func (e *EnqueueRequestForController) Update(evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	if req, err := e.getOwnerReconcileRequest(evt.MetaOld); err == nil {
		q.Add(req)
	}
	if req, err := e.getOwnerReconcileRequest(evt.MetaNew); err == nil {
		q.Add(req)
	}
}

// Delete implements EventHandler
func (e *EnqueueRequestForController) Delete(evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if req, err := e.getOwnerReconcileRequest(evt.Meta); err == nil {
		q.Add(req)
	}
}

// Generic implements EventHandler
func (e *EnqueueRequestForController) Generic(evt event.GenericEvent, q workqueue.RateLimitingInterface) {
	if req, err := e.getOwnerReconcileRequest(evt.Meta); err == nil {
		q.Add(req)
	}
}

// parseOwnerTypeGroupKind parses the OwnerType into a Group and Kind and caches the result.  Returns false
// if the OwnerType could not be parsed using the scheme.
func (e *EnqueueRequestForController) parseOwnerTypeGroupKind(scheme *runtime.Scheme) error {
	// Get the kinds of the type
	e.Scheme = scheme
	kinds, _, err := scheme.ObjectKinds(e.OwnerType)
	if err != nil {
		log.Error(err, "Could not get ObjectKinds for OwnerType", "owner type", fmt.Sprintf("%T", e.OwnerType))
		return err
	}
	// Expect only 1 kind.  If there is more than one kind this is probably an edge case such as ListOptions.
	if len(kinds) != 1 {
		err := fmt.Errorf("Expected exactly 1 kind for OwnerType %T, but found %s kinds", e.OwnerType, kinds)
		log.Error(nil, "Expected exactly 1 kind for OwnerType", "owner type", fmt.Sprintf("%T", e.OwnerType), "kinds", kinds)
		return err

	}
	// Cache the Group and Kind for the OwnerType
	e.groupKind = schema.GroupKind{Group: kinds[0].Group, Kind: kinds[0].Kind}
	return nil
}

// getOwnerReconcileRequest looks at object and returns a slice of reconcile.Request to reconcile
// owners of object that match e.OwnerType.
func (e *EnqueueRequestForController) getOwnerReconcileRequest(object metav1.Object) (reconcile.Request, error) {
	ref, err := e.getOwnerReference(object)
	if err != nil {
		owner, err := e.GetOwner(object)
		if err != nil {
			return reconcile.Request{}, err
		}

		return reconcile.Request{
			NamespacedName: types.NamespacedName{
				Namespace: owner.GetNamespace(),
				Name:      owner.GetName(),
			},
		}, nil
	}

	return reconcile.Request{
		NamespacedName: types.NamespacedName{
			Namespace: object.GetNamespace(),
			Name:      ref.Name,
		},
	}, nil
}

// getOwnersReferences returns the OwnerReferences for an object as specified by the EnqueueRequestForController
// - if IsController is true: only take the Controller OwnerReference (if found)
// - if IsController is false: take all OwnerReferences
// - if no owners are found try to get the object adopted
func (e *EnqueueRequestForController) getOwnerReference(object metav1.Object) (metav1.OwnerReference, error) {
	if object == nil {
		return metav1.OwnerReference{}, fmt.Errorf("cannot get owner reference of null object")
	}

	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If filtered to a Controller, only take the Controller OwnerReference
		return *ownerRef, nil
	}

	return metav1.OwnerReference{}, fmt.Errorf("no owner ref for object")
}

var _ inject.Scheme = &EnqueueRequestForController{}

// InjectScheme is called by the Controller to provide a singleton scheme to the EnqueueRequestForController.
func (e *EnqueueRequestForController) InjectScheme(s *runtime.Scheme) error {
	return e.parseOwnerTypeGroupKind(s)
}
