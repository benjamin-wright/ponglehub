package main_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type testReceiver struct {
	events chan event.Event
	cancel context.CancelFunc
	users  *client.UserClient
}

func newTestReceiver(u *testing.T, url string, subject string) *testReceiver {
	p, err := cenats.NewConsumer(url, subject, cenats.NatsOptions())
	if err != nil {
		assert.FailNow(u, "Error creating test receiver: %+v", err)
	}

	events, err := cloudevents.NewClient(p)
	if err != nil {
		assert.FailNow(u, "error creating cloudevents instance: %+v", err)
	}

	received := make(chan event.Event)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		err = events.StartReceiver(ctx, func(e event.Event) {
			received <- e
		})

		if err != nil {
			assert.FailNow(u, "Failed to start received: %+v", err)
		}
	}()

	users, err := client.New()
	if err != nil {
		assert.FailNow(u, "Failed to start client: %+v", err)
	}

	return &testReceiver{
		events: received,
		cancel: cancel,
		users:  users,
	}
}

func (t *testReceiver) deleteIfExists(u *testing.T, name string) {
	_, err := t.users.Get(name, v1.GetOptions{})
	if errors.IsNotFound(err) {
		return
	} else if err != nil {
		u.Errorf("Error getting test user: %+v", err)
		u.FailNow()
	}

	err = t.users.Delete(name, v1.DeleteOptions{})
	if err != nil {
		u.Errorf("Failed to delete test user: %+v", err)
		u.FailNow()
	}

	<-t.events
}

func (t *testReceiver) addUser(u *testing.T, name string, email string, password string) {
	err := t.users.Create(client.AuthUser{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Spec: client.AuthUserSpec{
			Name:     name,
			Email:    email,
			Password: password,
		},
	}, v1.CreateOptions{})

	if err != nil {
		assert.FailNow(u, "failed to create user: %+v", err)
	}
}

func (t *testReceiver) deleteUser(u *testing.T, name string) {
	err := t.users.Delete(name, v1.DeleteOptions{})

	if err != nil {
		assert.FailNow(u, "failed to delete user: %+v", err)
	}
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
