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

	foogroupv1 "github.com/raker22/k8s-foo-resource/pkg/apis/foogroup/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	err = c.Watch(&source.Kind{Type: &foogroupv1.FooReplicaSet{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by FooReplicaSet - change this for objects you create
	err = c.Watch(&source.Kind{Type: &foogroupv1.Foo{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &foogroupv1.FooReplicaSet{},
	})
	if err != nil {
		return err
	}

	return nil
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
// +kubebuilder:rbac:groups=foogroup.raker22.com,resources=foos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foogroup.raker22.com,resources=foos/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=foogroup.raker22.com,resources=fooreplicasets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foogroup.raker22.com,resources=fooreplicasets/status,verbs=get;update;patch
func (r *ReconcileFooReplicaSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the FooReplicaSet instance
	instance := &foogroupv1.FooReplicaSet{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if reflect.DeepEqual(instance.Spec.Selector, metav1.LabelSelector{}) {
		log.Info("Foo replica set is selecting all foos", "FooReplicaSet", instance.Name)
		return reconcile.Result{}, nil
	}

	selector, err := metav1.LabelSelectorAsSelector(&instance.Spec.Selector)

	fooList := &foogroupv1.FooList{}
	if err := r.List(context.TODO(), &client.ListOptions{
		LabelSelector: selector,
	}, fooList); err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this to be the object type created by your controller
	// Define the desired Foo object
	fooTemplate := &foogroupv1.Foo{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    instance.Namespace,
			GenerateName: fmt.Sprintf("%s-", instance.Name),
		},
		Spec: *instance.Spec.Template.DeepCopy(),
	}

	diffFoos := instance.Spec.Replicas - len(fooList.Items)

	if diffFoos > 0 {
		for i := 0; i < diffFoos; i++ {
			foo := fooTemplate.DeepCopy()

			if err := controllerutil.SetControllerReference(instance, foo, r.scheme); err != nil {
				return reconcile.Result{}, err
			}

			// TODO(user): Change this for the object type created by your controller
			log.Info("Creating Foo", "namespace", foo.Namespace, "name", foo.Name)
			if err := r.Create(context.TODO(), foo); err != nil {
				log.Info("Failed to create Foo", "namespace", foo.Namespace, "FooReplicaSet", instance.Name)
				return reconcile.Result{}, err
			}
		}
	}

	log.Info("Foos matching foo replica set", "FooReplicaSet", instance.Name)
	for _, foo := range fooList.Items {
		log.Info("", "Foo", foo.Name)
	}

	//
	//// TODO(user): Change this for the object type created by your controller
	//// Update the found object and write the result back if there are any changes
	//if !reflect.DeepEqual(foo.Spec, found.Spec) {
	//	found.Spec = foo.Spec
	//	log.Info("Updating Foo", "namespace", foo.Namespace, "name", foo.Name)
	//	//err = r.Update(context.TODO(), found)
	//	if err != nil {
	//		return reconcile.Result{}, err
	//	}
	//}

	return reconcile.Result{}, nil
}
