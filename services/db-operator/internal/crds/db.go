package crds

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type CockroachDBStatus struct {
	Ready bool `json:"ready"`
}

type CockroachDBSpec struct {
	Storage string `json:"storage"`
}

type CockroachDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CockroachDBSpec   `json:"spec"`
	Status CockroachDBStatus `json:"status"`
}

type CockroachDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CockroachDB `json:"items"`
}
