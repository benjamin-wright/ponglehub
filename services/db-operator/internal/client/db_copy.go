package client

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *CockroachDB) DeepCopyInto(out *CockroachDB) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = CockroachDBSpec{
		Storage: in.Spec.Storage,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *CockroachDB) DeepCopyObject() runtime.Object {
	out := CockroachDB{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *CockroachDBList) DeepCopyObject() runtime.Object {
	out := CockroachDBList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]CockroachDB, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
