# Types

Add properties to `Spec`/`Status`.
Resource and List types are automatically generated with the following template.

```go
// where `Kind` is the resource kind from `kubebuilder create api`
type Kind struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KindSpec   `json:"spec,omitempty"`
	Status KindStatus `json:"status,omitempty"`
}

type KindList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kind `json:"items"`
}
```

`make install` generates the [CRDs](../../../../config/crds) in `config/crds` and applies them to the cluster.
It should be run any time a type is updated to ensure the CRDs are up to date.

```go
// foo_types.go
// FooSpec defines the desired state of Foo
type FooSpec struct {
	Message string `json:"message"`
	Value   int32  `json:"value"`
}

// FooStatus defines the observed state of Foo
type FooStatus struct {
	Message string `json:"message"`
	Value   int32  `json:"value"`
}
```

An example of validation on the `Replicas` property.

```go
// foorepicaset_types.go
// FooReplicaSetSpec defines the desired state of FooReplicaSet
type FooReplicaSetSpec struct {
	// +kubebuilder:validation:Minimum=1
	Replicas int                  `json:"replicas"`
	Template Foo                  `json:"template"`
	Selector metav1.LabelSelector `json:"selector"`
}

// FooReplicaSetStatus defines the observed state of FooReplicaSet
// could generate without status subresource
type FooReplicaSetStatus struct {
	
}
```


