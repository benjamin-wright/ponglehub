package crds

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *EventTrigger) DeepCopyInto(out *EventTrigger) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta

	out.Spec = EventTriggerSpec{
		Filters: make([]string, len(in.Spec.Filters)),
		URL:     in.Spec.URL,
	}

	copy(out.Spec.Filters, in.Spec.Filters)
}

// DeepCopyObject returns a generically typed copy of an object
func (in *EventTrigger) DeepCopyObject() runtime.Object {
	out := EventTrigger{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *EventTriggerList) DeepCopyObject() runtime.Object {
	out := EventTriggerList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]EventTrigger, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
