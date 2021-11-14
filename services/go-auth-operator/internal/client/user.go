package client

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type AuthUserSpec struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUserStatus struct {
	ID      string `json:"id"`
	Pending bool   `json:"pending"`
}

type AuthUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AuthUserSpec   `json:"spec"`
	Status AuthUserStatus `json:"status"`
}

func (a *AuthUser) Equals(user *AuthUser) bool {
	return a.Spec.Email == user.Spec.Email &&
		a.Spec.Name == user.Spec.Name &&
		a.Spec.Password == user.Spec.Password
}

type AuthUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AuthUser `json:"items"`
}
