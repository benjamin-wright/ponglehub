package integration

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/auth/auth-operator/internal/client"
)

func getEnv(t *testing.T, key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		t.Errorf("Missing environment variable %s", key)
		t.FailNow()
	}

	return value
}

func TestCRDCrud(t *testing.T) {
	client.AddToScheme(scheme.Scheme)

	natsUrl := getEnv(t, "NATS_URL")
	natsSubject := getEnv(t, "NATS_SUBJECT")

	for _, test := range []struct {
		Name         string
		Prepare      func(*testReceiver, *testing.T)
		Send         func(*testReceiver, *testing.T)
		ExpectedType string
		ExpectedData interface{}
	}{
		{
			Name: "Add user",
			Prepare: func(r *testReceiver, u *testing.T) {
				r.deleteIfExists(u, "test-user")
			},
			Send: func(r *testReceiver, u *testing.T) {
				r.addUser(u, "test-user", "test@user.com", "P@ssw0rd123!")
			},
			ExpectedType: "ponglehub.auth.user.add",
			ExpectedData: map[string]string{"meta_name": "test-user", "name": "test-user", "email": "test@user.com", "password": "P@ssw0rd123!"},
		},
		{
			Name: "Delete user",
			Prepare: func(r *testReceiver, u *testing.T) {
				r.deleteIfExists(u, "test-user")
				r.addUser(u, "test-user", "test@user.com", "P@ssw0rd123!")
				<-r.events
			},
			Send: func(r *testReceiver, u *testing.T) {
				r.deleteUser(u, "test-user")
			},
			ExpectedType: "ponglehub.auth.user.delete",
			ExpectedData: map[string]string{"meta_name": "test-user", "name": "test-user"},
		},
		{
			Name: "Update user",
			Prepare: func(r *testReceiver, u *testing.T) {
				r.deleteIfExists(u, "test-user")
				r.addUser(u, "test-user", "test@user.com", "P@ssw0rd123!")
				<-r.events
			},
			Send: func(r *testReceiver, u *testing.T) {
				r.updateUser(u, "test-user", "new@email.com", "newP@ssw0rd")
			},
			ExpectedType: "ponglehub.auth.user.update",
			ExpectedData: map[string]string{"meta_name": "test-user", "name": "test-user", "email": "new@email.com", "password": "newP@ssw0rd"},
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			receiver := newTestReceiver(u, natsUrl, natsSubject)
			defer receiver.cancel()

			test.Prepare(receiver, u)
			test.Send(receiver, u)

			e := <-receiver.events

			var actual map[string]string
			err := json.Unmarshal(e.Data(), &actual)
			assert.NoError(u, err)
			assert.Equal(u, test.ExpectedType, e.Type())
			assert.Equal(u, test.ExpectedData, actual)
		})
	}
}
