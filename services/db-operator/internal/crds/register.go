package crds

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "ponglehub.co.uk"
const GroupVersion = "v1alpha1"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&CockroachClient{},
		&CockroachClientList{},
	)

	scheme.AddKnownTypes(SchemeGroupVersion,
		&CockroachDB{},
		&CockroachDBList{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}