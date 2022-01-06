package crds

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Use this object outside of the package
type User struct {
	ID              string
	Name            string
	ResourceVersion string
	Display         string
	Email           string
	Invited         bool
	Member          bool
}

type AuthUserSpec struct {
	Display string `json:"display"`
	Email   string `json:"email"`
}

type AuthUserStatus struct {
	Invited bool `json:"invited"`
	Member  bool `json:"member"`
}

// For internal use only
type AuthUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuthUserSpec   `json:"spec"`
	Status AuthUserStatus `json:"status"`
}

type AuthUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AuthUser `json:"items"`
}
