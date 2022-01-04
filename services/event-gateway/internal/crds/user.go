package crds

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Use this object outside of the package
type User struct {
	ID              string
	Name            string
	ResourceVersion string
	Display         string
	Email           string
	InviteToken     string
	PasswordHash    string
}

type AuthUserSpec struct {
	Display      string `json:"display"`
	Email        string `json:"email"`
	InviteToken  string `json:"inviteToken"`
	PasswordHash string `json:"passwordHash"`
}

// For internal use only
type AuthUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec AuthUserSpec `json:"spec"`
}

type AuthUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AuthUser `json:"items"`
}
