/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fooreplicaset

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/types"

	foov1 "github.com/raker22/k8s-foo-resource/pkg/apis/foo/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new FooReplicaSet Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileFooReplicaSet{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("fooreplicaset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to FooReplicaSet
	err = c.Watch(&source.Kind{Type: &foov1.FooReplicaSet{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by FooReplicaSet - change this for objects you create
	err = c.Watch(&source.Kind{Type: &foov1.Foo{}}, &EnqueueRequestForOrphans{
		EnqueueRequestForOwner: handler.EnqueueRequestForOwner{
			OwnerType:    &foov1.FooReplicaSet{},
			IsController: true,
		},
		Adopt: func(e *EnqueueRequestForOrphans, object metav1.Object) ([]metav1.OwnerReference, error) {
			log.Info("Attempting to adopt foo", "foo", object.GetName())
			namespacedName := types.NamespacedName{
				Namespace: object.GetNamespace(),
				Name:      object.GetName(),
			}

			rec := ReconcileFooReplicaSet(r)

			foo := &foov1.Foo{}
			if err := rec.Get(context.TODO(), namespacedName, foo); err != nil {
				return []metav1.OwnerReference{}, err
			}

			fooReplicaSetList := &foov1.FooReplicaSetList{}
			if err := rec.List(context.TODO(), &client.ListOptions{}, fooReplicaSetList); err != nil {
				return []metav1.OwnerReference{}, err
			}

			fooLabels := labels.Set(foo.Labels)

			for _, fooReplicaSet := range fooReplicaSetList.Items {
				if reflect.DeepEqual(fooReplicaSet.Spec.Selector, metav1.LabelSelector{}) {
					log.Info("Foo replica set is selecting all foos", "fooReplicaSet", fooReplicaSet.Name)
					continue
				}

				selector, err := metav1.LabelSelectorAsSelector(&fooReplicaSet.Spec.Selector)
				if err != nil {
					return []metav1.OwnerReference{}, nil
				}

				if selector.Matches(fooLabels) {
					log.Info("Adopting foo",
						"namespace", foo.Namespace,
						"fooReplicaSet", fooReplicaSet.Name,
						"foo", foo.Name,
					)
					if err := controllerutil.SetControllerReference(&fooReplicaSet, foo, r.scheme); err != nil {
						log.Info("Failed to adopt foo",
							"namespace", foo.Namespace,
							"fooReplicaSet", fooReplicaSet.Name,
							"foo", foo.Name,
						)
					} else {
						if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
							return []metav1.OwnerReference{*ownerRef}, nil
						}
					}
				}
			}

			return []metav1.OwnerReference{}, fmt.Errorf("no possible owners for object")
		},
	})
	if err != nil {
		return err
	}

	return nil
}

type EnqueueRequestForOrphans struct {
	handler.EnqueueRequestForOwner
	Adopt func(e *EnqueueRequestForOrphans, object metav1.Object) ([]metav1.OwnerReference, error)
}

func (e *EnqueueRequestForOrphans) getOwnersReferences(object metav1.Object) []metav1.OwnerReference {
	if object == nil {
		return nil
	}

	var refs []metav1.OwnerReference

	// If not filtered as Controller only, then use all the OwnerReferences
	if !e.IsController {
		refs = object.GetOwnerReferences()
	}
	// If filtered to a Controller, only take the Controller OwnerReference
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		refs = []metav1.OwnerReference{*ownerRef}
	}

	if len(refs) == 0 {
		ownerRefs, err := e.Adopt(e, object)
		if err == nil {
			refs = ownerRefs
		}
	}

	return refs
}

var _ reconcile.Reconciler = &ReconcileFooReplicaSet{}

// ReconcileFooReplicaSet reconciles a FooReplicaSet object
type ReconcileFooReplicaSet struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a FooReplicaSet object and makes changes based on the state read
// and what is in the FooReplicaSet.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=foo.raker22.com,resources=foos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foo.raker22.com,resources=foos/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=foo.raker22.com,resources=fooreplicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foo.raker22.com,resources=fooreplicasets/status,verbs=get;update;patch
func (r *ReconcileFooReplicaSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the FooReplicaSet instance
	log.Info("reconciling foo replica set", "namespace", request.NamespacedName.Namespace, "name", request.NamespacedName.Name)
	instance := &foov1.FooReplicaSet{}
	if err := r.Get(context.TODO(), request.NamespacedName, instance); err != nil {
		// check if we got a foo instead of a foo replica set
		fooInstance := &foov1.Foo{}
		if er := r.Get(context.TODO(), request.NamespacedName, fooInstance); er != nil {
			// Error reading the object
			if errors.IsNotFound(err) {
				// Object not found, return.  Created objects are automatically garbage collected.
				// For additional cleanup logic use finalizers.
				return reconcile.Result{}, nil
			}

			// requeue the request.
			return reconcile.Result{}, err
		}

		// it was a foo so get matching replica sets
		log.Info("it was a foo", "namespace", request.NamespacedName.Namespace, "name", request.NamespacedName.Name)
		return reconcile.Result{}, nil
	}

	if reflect.DeepEqual(instance.Spec.Selector, metav1.LabelSelector{}) {
		log.Info("Foo replica set is selecting all foos", "fooReplicaSet", instance.Name)
		return reconcile.Result{}, nil
	}

	selector, err := metav1.LabelSelectorAsSelector(&instance.Spec.Selector)
	if err != nil {
		return reconcile.Result{}, err
	}

	if !selector.Matches(labels.Set(instance.Spec.Template.Labels)) {
		log.Info("Invalid foo replica set: selector must match template labels",
			"namespace", instance.Namespace,
			"fooReplicaSet", instance.Name,
			"selector", instance.Spec.Selector,
			"labels", instance.Spec.Template.Labels,
		)
		return reconcile.Result{}, nil
	}

	fooList := &foov1.FooList{}
	if err := r.List(context.TODO(), &client.ListOptions{
		LabelSelector: selector,
	}, fooList); err != nil {
		return reconcile.Result{}, err
	}

	//if instance.DeletionTimestamp != nil {
	//	log.Info("Foo replica set deleted", "namespace", instance.Namespace, "fooReplicaSet", instance.Name)
	//	for _, foo := range fooList.Items {
	//		log.Info("Deleting foo", "namespace", foo.Namespace, "foo", foo.Name)
	//		if err := r.Delete(context.TODO(), &foo); err != nil {
	//			return reconcile.Result{}, err
	//		}
	//	}
	//	return reconcile.Result{}, nil
	//}

	// TODO(user): Change this to be the object type created by your controller
	// Define the desired Foo object
	fooTemplate := instance.Spec.Template.DeepCopy()
	fooTemplate.Status = foov1.FooStatus{}
	fooTemplate.Namespace = instance.Namespace
	fooTemplate.GenerateName = fmt.Sprintf("%s-", instance.Name)

	//gvk, err := apiutil.GVKForObject(instance, r.scheme)
	//if err != nil {
	//	return reconcile.Result{}, err
	//}

	var updateFoos []foov1.Foo
	for _, foo := range fooList.Items {
		ownerRef := metav1.GetControllerOf(&foo)
		//hasOwner := false
		//isOwner := false
		//
		//for _, ref := range ownerRefs {
		//	if ref.Kind == gvk.Kind {
		//		hasOwner = true
		//
		//		if ref.UID == instance.UID {
		//			isOwner = true
		//			break
		//		}
		//	}
		//}
		//
		//if !hasOwner {
		//	log.Info("Adopting foo",
		//		"namespace", instance.Namespace,
		//		"fooReplicaSet", instance.Name,
		//		"foo", foo.Name,
		//	)
		//	if err := controllerutil.SetControllerReference(instance, &foo, r.scheme); err != nil {
		//		log.Info("Failed to adopt foo",
		//			"namespace", instance.Namespace,
		//			"fooReplicaSet", instance.Name,
		//			"foo", foo.Name,
		//		)
		//		return reconcile.Result{}, err
		//	}
		//	isOwner = true
		//}

		if ownerRef.UID == instance.UID {
			updateFoos = append(updateFoos, foo)
		}
	}
	diffFoos := instance.Spec.Replicas - len(updateFoos)

	if diffFoos > 0 {
		for i := 0; i < diffFoos; i++ {
			foo := fooTemplate.DeepCopy()

			if err := controllerutil.SetControllerReference(instance, foo, r.scheme); err != nil {
				log.Info("Failed to set owner of foo",
					"namespace", instance.Namespace,
					"fooReplicaSet", instance.Name,
				)
				return reconcile.Result{}, err
			}

			// TODO(user): Change this for the object type created by your controller
			log.Info("Creating foo", "namespace", foo.Namespace, "generateName", foo.GenerateName)
			if err := r.Create(context.TODO(), foo); err != nil {
				log.Info("Failed to create foo", "namespace", foo.Namespace, "fooReplicaSet", instance.Name)
				return reconcile.Result{}, err
			}
		}
	} else if diffFoos < 0 {
		i := -diffFoos
		removeFoos := fooList.Items[:i]
		updateFoos = fooList.Items[i:]

		for _, foo := range removeFoos {
			log.Info("Deleting foo",
				"namespace", foo.Namespace,
				"fooReplicaSet", instance.Name,
				"foo", foo.Name,
			)
			if err := r.Delete(context.TODO(), &foo); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	for _, foo := range updateFoos {
		// TODO(user): Change this for the object type created by your controller
		// Update the found object and write the result back if there are any changes
		if !reflect.DeepEqual(fooTemplate.Spec, foo.Spec) {
			foo.Spec = fooTemplate.Spec
			log.Info("Updating foo",
				"namespace", foo.Namespace,
				"fooReplicaSet", instance.Name,
				"foo", foo.Name,
			)
			err = r.Update(context.TODO(), &foo)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	return reconcile.Result{}, nil
}
