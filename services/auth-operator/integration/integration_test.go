package integration

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/auth/auth-operator/internal/users"
	"ponglehub.co.uk/events/recorder/pkg/recorder"
	events "ponglehub.co.uk/lib/user-events"
)

func setNoUser(u *testing.T, userClient *users.UserClient) {
	userClient.Delete("test-user")
}

func setUser(u *testing.T, userClient *users.UserClient, user events.User) {
	userClient.Delete("test-user")

	_, err := userClient.Create(events.User{
		Name:     "test-user",
		Username: "test-user",
		Email:    "test@user.com",
		Password: "P@ssw0rd123!",
	})
	if err != nil {
		assert.FailNow(u, "Failed to create test user %+v", err)
	}

	time.Sleep(250 * time.Millisecond)

	current, err := userClient.Get("test-user")
	if err != nil {
		assert.FailNow(u, "Failed to get test user %+v", err)
	}

	user.ResourceVersion = current.ResourceVersion
	_, err = userClient.Status(user)
	if err != nil {
		assert.FailNow(u, "Failed to get test user %+v", err)
	}
}

func TestCRDCrud(t *testing.T) {
	RECORDER_URL := os.Getenv("RECORDER_URL")

	users.AddToScheme(scheme.Scheme)
	userClient, err := users.New(&users.ClientArgs{External: true})
	if err != nil {
		assert.FailNow(t, "failed to start users client: %+v", err)
	}

	eventClient, err := events.New("OPERATOR_URL", "test-operator")
	if err != nil {
		assert.FailNow(t, "failed to start event client: %+v", err)
	}

	for _, test := range []struct {
		Name         string
		Prepare      func(*testing.T)
		Send         func(*testing.T)
		ExpectedType string
		ExpectedData events.User
		ExpectUser   bool
	}{
		{
			Name: "Add user",
			Prepare: func(u *testing.T) {
				setNoUser(u, userClient)
			},
			Send: func(u *testing.T) {
				_, err := userClient.Create(events.User{
					Name:     "test-user",
					Username: "test-user",
					Email:    "test@user.com",
					Password: "P@ssw0rd123!",
				})
				assert.NoError(u, err)
			},
			ExpectedType: "ponglehub.auth.user.add",
			ExpectedData: events.User{
				Name:     "test-user",
				Username: "test-user",
				Email:    "test@user.com",
				Password: "P@ssw0rd123!",
				Pending:  true,
				ID:       "",
			},
			ExpectUser: true,
		},
		{
			Name: "Delete user",
			Prepare: func(u *testing.T) {
				setUser(u, userClient, events.User{
					Name:     "test-user",
					Username: "test-user",
					Email:    "test@user.com",
					Password: "P@ssw0rd123!",
					Pending:  true,
				})
			},
			Send: func(u *testing.T) {
				assert.NoError(u, userClient.Delete("test-user"))
			},
			ExpectedType: "ponglehub.auth.user.delete",
			ExpectedData: events.User{
				Name:     "test-user",
				Username: "test-user",
				Email:    "test@user.com",
				Password: "P@ssw0rd123!",
				Pending:  true,
				ID:       "",
			},
		},
		{
			Name: "Update user",
			Prepare: func(u *testing.T) {
				setUser(u, userClient, events.User{
					Name:     "test-user",
					Username: "test-user",
					Email:    "test@user.com",
					Password: "P@ssw0rd123!",
				})
			},
			Send: func(u *testing.T) {
				user, err := userClient.Get("test-user")
				if err != nil {
					assert.FailNow(u, "Failed to get test user: %+v", err)
				}

				user.Email = "new@email.com"
				user.Password = "newP@ssw0rd"

				_, err = userClient.Update(user)
				assert.NoError(u, err)
			},
			ExpectedType: "ponglehub.auth.user.update",
			ExpectedData: events.User{
				Name:     "test-user",
				Username: "test-user",
				Email:    "new@email.com",
				Password: "newP@ssw0rd",
				Pending:  true,
				ID:       "",
			},
			ExpectUser: true,
		},
		{
			Name: "Set user event",
			Prepare: func(u *testing.T) {
				setUser(u, userClient, events.User{
					Name:     "test-user",
					Username: "test-user",
					Email:    "test@user.com",
					Password: "P@ssw0rd123!",
				})
			},
			Send: func(u *testing.T) {
				user, err := userClient.Get("test-user")
				if err != nil {
					assert.FailNow(u, "Failed to get test user: %+v", err)
				}

				user.ID = "1234"
				eventClient.SetUser(user)
			},
			ExpectedType: "ponglehub.auth.user.set.ack",
			ExpectedData: events.User{
				Name:     "test-user",
				Username: "test-user",
				Email:    "test@user.com",
				Password: "P@ssw0rd123!",
				Pending:  false,
				ID:       "1234",
			},
			ExpectUser: true,
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			test.Prepare(u)
			recorder.Clear(u, RECORDER_URL)
			test.Send(u)

			data := recorder.WaitForEvent(u, RECORDER_URL, test.ExpectedType)
			eventData := events.User{}
			assert.NoError(u, json.Unmarshal([]byte(data), &eventData))

			eventData.ResourceVersion = ""

			assert.NoError(u, err)
			assert.Equal(u, test.ExpectedData, eventData)

			if test.ExpectUser {
				user, err := userClient.Get(test.ExpectedData.Name)
				assert.NoError(u, err)

				user.ResourceVersion = ""

				assert.Equal(u, test.ExpectedData, user)
			}
		})
	}
}
