package integration

import (
	"context"
	"testing"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ponglehub.co.uk/auth/auth-operator/internal/client"
)

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
	err := t.users.Create(&client.AuthUser{
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

func (t *testReceiver) updateUser(u *testing.T, name string, email string, password string) {
	user, err := t.users.Get(name, v1.GetOptions{})
	if err != nil {
		assert.FailNow(u, "failed to get user for status update: %+v", err)
	}

	user.Status.Pending = false
	err = t.users.SetStatus(user, v1.UpdateOptions{})
	if err != nil {
		assert.FailNow(u, "failed to reset user pending status: %+v", err)
	}

	user, err = t.users.Get(name, v1.GetOptions{})
	if err != nil {
		assert.FailNow(u, "failed to get user for value update: %+v", err)
	}

	user.Spec.Email = email
	user.Spec.Password = password

	err = t.users.Update(user, v1.UpdateOptions{})

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
