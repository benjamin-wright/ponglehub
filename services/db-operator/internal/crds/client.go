package crds

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type CockroachClientSpec struct {
	Database string `json:"database"`
	Secret   string `json:"secret"`
}

type CockroachClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CockroachClientSpec `json:"spec"`
}

type CockroachClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CockroachClient `json:"items"`
}
