package crds

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type CockroachClientSpec struct {
	Deployment string `json:"deployment"`
	Database   string `json:"database"`
	Username   string `json:"username"`
	Secret     string `json:"secret"`
}

type CockroachClientStatus struct {
	Ready bool `json:"ready"`
}

type CockroachClient struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CockroachClientSpec   `json:"spec"`
	Status CockroachClientStatus `json:"status"`
}

type CockroachClientList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []CockroachClient `json:"items"`
}
