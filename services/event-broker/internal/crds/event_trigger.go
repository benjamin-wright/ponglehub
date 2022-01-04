package crds

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type EventTriggerSpec struct {
	URL     string   `json:"url"`
	Filters []string `json:"filters"`
}

type EventTrigger struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec EventTriggerSpec `json:"spec"`
}

type EventTriggerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []EventTrigger `json:"items"`
}
