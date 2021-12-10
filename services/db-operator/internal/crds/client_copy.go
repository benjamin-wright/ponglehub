package crds

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *CockroachClient) DeepCopyInto(out *CockroachClient) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = CockroachClientSpec{
		Database: in.Spec.Database,
		Secret:   in.Spec.Secret,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *CockroachClient) DeepCopyObject() runtime.Object {
	out := CockroachClient{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *CockroachClientList) DeepCopyObject() runtime.Object {
	out := CockroachClientList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]CockroachClient, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
