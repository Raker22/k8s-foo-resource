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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type FooTemplate struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              FooSpec `json:"spec"`
}

// FooReplicaSetSpec defines the desired state of FooReplicaSet
type FooReplicaSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +kubebuilder:validation:Minimum=1
	Replicas int                  `json:"replicas"`
	Template FooTemplate          `json:"template"`
	Selector metav1.LabelSelector `json:"selector"`
}

// FooReplicaSetStatus defines the observed state of FooReplicaSet
type FooReplicaSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	CurrentReplicas int `json:"currentReplicas"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooReplicaSet is the Schema for the fooreplicasets API
// +k8s:openapi-gen=true
// +kubebuilder:categories=all,allfoo
// +kubebuilder:printcolumn:name="Current Replicas",type="integer",JSONPath=".status.currentReplicas"
// +kubebuilder:printcolumn:name="Desired Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Message",type="string",JSONPath=".spec.template.spec.message"
// +kubebuilder:printcolumn:name="Value",type="integer",JSONPath=".spec.template.spec.value"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type FooReplicaSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooReplicaSetSpec   `json:"spec,omitempty"`
	Status FooReplicaSetStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooReplicaSetList contains a list of FooReplicaSet
type FooReplicaSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FooReplicaSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FooReplicaSet{}, &FooReplicaSetList{})
}
