package client

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *AuthUser) DeepCopyInto(out *AuthUser) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = AuthUserSpec{
		Name:     in.Spec.Name,
		Email:    in.Spec.Email,
		Password: in.Spec.Password,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *AuthUser) DeepCopyObject() runtime.Object {
	out := AuthUser{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *AuthUserList) DeepCopyObject() runtime.Object {
	out := AuthUserList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]AuthUser, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
