package integration

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/auth/auth-operator/internal/client"
)

func TestCRDCrud(t *testing.T) {
	client.AddToScheme(scheme.Scheme)
	receiver := newTestReceiver(t)
	defer receiver.cancel()

	for _, test := range []struct {
		Name         string
		Prepare      func(*testing.T)
		Send         func(*testing.T)
		ExpectedType string
		ExpectedData interface{}
	}{
		{
			Name: "Add user",
			Prepare: func(u *testing.T) {
				receiver.deleteIfExists(u, "test-user")
			},
			Send: func(u *testing.T) {
				receiver.addUser(u, "test-user", "test@user.com", "P@ssw0rd123!")
			},
			ExpectedType: "ponglehub.auth.user.add",
			ExpectedData: map[string]string{"meta_name": "test-user", "name": "test-user", "email": "test@user.com", "password": "P@ssw0rd123!"},
		},
		{
			Name: "Delete user",
			Prepare: func(u *testing.T) {
				receiver.deleteIfExists(u, "test-user")
				receiver.addUser(u, "test-user", "test@user.com", "P@ssw0rd123!")
				<-receiver.events
			},
			Send: func(u *testing.T) {
				receiver.deleteUser(u, "test-user")
			},
			ExpectedType: "ponglehub.auth.user.delete",
			ExpectedData: map[string]string{"meta_name": "test-user", "name": "test-user"},
		},
		{
			Name: "Update user",
			Prepare: func(u *testing.T) {
				receiver.deleteIfExists(u, "test-user")
				receiver.addUser(u, "test-user", "test@user.com", "P@ssw0rd123!")
				<-receiver.events
			},
			Send: func(u *testing.T) {
				receiver.updateUser(u, "test-user", "new@email.com", "newP@ssw0rd")
			},
			ExpectedType: "ponglehub.auth.user.update",
			ExpectedData: map[string]string{"meta_name": "test-user", "name": "test-user", "email": "new@email.com", "password": "newP@ssw0rd"},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			test.Prepare(u)
			test.Send(u)

			e := <-receiver.events

			var actual map[string]string
			err := json.Unmarshal(e.Data(), &actual)
			assert.NoError(u, err)
			assert.Equal(u, test.ExpectedType, e.Type())
			assert.Equal(u, test.ExpectedData, actual)
		})
	}
}
