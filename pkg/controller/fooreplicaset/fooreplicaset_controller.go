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
	"reflect"

	foov1 "github.com/raker22/k8s-foo-resource/pkg/apis/foo/v1"
	foohandler "github.com/raker22/k8s-foo-resource/pkg/handler"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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

	// watch for changes to Foo that are owned by or could be owned by a FooReplicaSet
	err = c.Watch(&source.Kind{Type: &foov1.Foo{}}, &foohandler.EnqueueRequestForController{
		OwnerType: &foov1.FooReplicaSet{},
		GetOwner: func(object metav1.Object) (metav1.Object, error) {
			namespacedName := types.NamespacedName{
				Namespace: object.GetNamespace(),
				Name:      object.GetName(),
			}

			c := mgr.GetClient()

			foo := &foov1.Foo{}
			if err := c.Get(context.TODO(), namespacedName, foo); err != nil {
				return nil, err
			}

			fooReplicaSetList := &foov1.FooReplicaSetList{}
			if err := c.List(context.TODO(), &client.ListOptions{}, fooReplicaSetList); err != nil {
				return nil, err
			}

			fooLabels := labels.Set(foo.Labels)

			for _, fooReplicaSet := range fooReplicaSetList.Items {
				if reflect.DeepEqual(fooReplicaSet.Spec.Selector, metav1.LabelSelector{}) {
					log.Info("Foo replica set is selecting all foos", "fooReplicaSet", fooReplicaSet.Name)
					continue
				}

				selector, err := metav1.LabelSelectorAsSelector(&fooReplicaSet.Spec.Selector)
				if err != nil {
					log.Info("Error parsing foo replica set selector", "fooReplicaSet", fooReplicaSet.Name)
					continue
				}

				if selector.Matches(fooLabels) {
					return &fooReplicaSet, nil
				}
			}

			return nil, fmt.Errorf("no possible owner for object")
		},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileFooReplicaSet{}

// ReconcileFooReplicaSet reconciles a FooReplicaSet object
type ReconcileFooReplicaSet struct {
	reconcile.Reconciler
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a FooReplicaSet object and makes changes based on the state read
// and what is in the FooReplicaSet.Spec
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=foo.raker22.com,resources=foos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foo.raker22.com,resources=foos/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=foo.raker22.com,resources=fooreplicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foo.raker22.com,resources=fooreplicasets/status,verbs=get;update;patch
func (r *ReconcileFooReplicaSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the FooReplicaSet instance
	instance := &foov1.FooReplicaSet{}
	if err := r.Get(context.TODO(), request.NamespacedName, instance); err != nil {
		// Error reading the object
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}

		// requeue the request.
		return reconcile.Result{}, err
	}

	if reflect.DeepEqual(instance.Spec.Selector, metav1.LabelSelector{}) {
		log.Info("Foo replica set is selecting all foos, skipping", "fooReplicaSet", instance.Name)
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
	if err := r.List(context.TODO(), &client.ListOptions{}, fooList); err != nil {
		return reconcile.Result{}, err
	}

	// Define the desired Foo object
	fooTemplate := instance.Spec.Template
	fooTemplate.Namespace = instance.Namespace
	fooTemplate.GenerateName = fmt.Sprintf("%s-", instance.Name)

	var ownFoos []foov1.Foo
	for _, foo := range fooList.Items {
		if selector.Matches(labels.Set(foo.Labels)) {
			// get own foos and adopt foos that don't have an owners
			if ownerRef := metav1.GetControllerOf(&foo); ownerRef != nil {
				if ownerRef.UID == instance.UID {
					ownFoos = append(ownFoos, foo)
				}
			} else {
				log.Info("Adopting foo",
					"namespace", instance.Namespace,
					"fooReplicaSet", instance.Name,
					"foo", foo.Name,
				)
				if err := controllerutil.SetControllerReference(instance, &foo, r.scheme); err == nil {
					ownFoos = append(ownFoos, foo)
				} else {
					log.Info("Failed to adopt foo",
						"namespace", instance.Namespace,
						"fooReplicaSet", instance.Name,
						"foo", foo.Name,
					)
				}
			}
		} else if ref := metav1.GetControllerOf(&foo); ref != nil && instance.UID == ref.UID {
			log.Info("Orphaning foo",
				"namespace", instance.Namespace,
				"fooReplicaSet", instance.Name,
				"foo", foo.Name,
			)
			var refs []metav1.OwnerReference
			for _, ref := range foo.OwnerReferences {
				if ref.UID != instance.UID {
					refs = append(refs, ref)
				}
			}

			foo.OwnerReferences = refs
			if err = r.Update(context.TODO(), &foo); err != nil {
				log.Info("Failed to orphan foo",
					"namespace", instance.Namespace,
					"fooReplicaSet", instance.Name,
					"foo", foo.Name,
				)
				return reconcile.Result{}, err
			}
		}
	}
	diffFoos := instance.Spec.Replicas - len(ownFoos)

	if diffFoos > 0 {
		for i := 0; i < diffFoos; i++ {
			foo := &foov1.Foo{
				ObjectMeta: fooTemplate.ObjectMeta,
				Spec:       fooTemplate.Spec,
			}

			if err := controllerutil.SetControllerReference(instance, foo, r.scheme); err != nil {
				log.Info("Failed to set owner of foo",
					"namespace", instance.Namespace,
					"fooReplicaSet", instance.Name,
				)
				return reconcile.Result{}, err
			}

			log.Info("Creating foo", "namespace", foo.Namespace, "generateName", foo.GenerateName)
			if err := r.Create(context.TODO(), foo); err != nil {
				log.Info("Failed to create foo", "namespace", foo.Namespace, "fooReplicaSet", instance.Name)
				return reconcile.Result{}, err
			}
		}
	} else if diffFoos < 0 {
		i := -diffFoos
		removeFoos := ownFoos[:i]
		ownFoos = ownFoos[i:]

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

	for _, foo := range ownFoos {
		// Update the found object and write the result back if there are any changes
		if !reflect.DeepEqual(fooTemplate.Spec, foo.Spec) {
			foo.Spec = fooTemplate.Spec
			log.Info("Updating foo",
				"namespace", foo.Namespace,
				"fooReplicaSet", instance.Name,
				"foo", foo.Name,
			)
			if err = r.Update(context.TODO(), &foo); err != nil {
				return reconcile.Result{}, err
			}
		}
	}

	status := instance.Status.DeepCopy()

	status.CurrentReplicas = instance.Spec.Replicas

	if !reflect.DeepEqual(*status, instance.Status) {
		instance.Status = *status
		log.Info("Updating foo replica set",
			"namespace", instance.Namespace,
			"fooReplicaSet", instance.Name,
		)
		if err = r.Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
