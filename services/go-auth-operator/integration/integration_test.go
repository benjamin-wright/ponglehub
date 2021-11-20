package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/auth/auth-operator/internal/events"
	"ponglehub.co.uk/auth/auth-operator/internal/users"
)

func TestCRDCrud(t *testing.T) {
	users.AddToScheme(scheme.Scheme)
	userClient, err := users.New()
	if err != nil {
		assert.FailNow(t, "failed to start users client: %+v", err)
	}

	receiver := make(chan events.UserEvent, 5)
	cancel, err := events.Listen(func(event events.UserEvent) {
		receiver <- event
	})
	if err != nil {
		assert.FailNow(t, "failed to start event listener: %+v", err)
	}
	defer cancel()

	for _, test := range []struct {
		Name         string
		Prepare      func(*testing.T)
		Send         func(*testing.T)
		ExpectedType string
		ExpectedData events.User
	}{
		{
			Name: "Add user",
			Prepare: func(u *testing.T) {
				if err := userClient.Delete("test-user"); err == nil {
					<-receiver
				}
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
		},
		{
			Name: "Delete user",
			Prepare: func(u *testing.T) {
				if err := userClient.Delete("test-user"); err == nil {
					<-receiver
				}
				userClient.Create(events.User{
					Name:     "test-user",
					Username: "test-user",
					Email:    "test@user.com",
					Password: "P@ssw0rd123!",
				})
				<-receiver
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
				if err := userClient.Delete("test-user"); err == nil {
					<-receiver
				}
				_, err := userClient.Create(events.User{
					Name:     "test-user",
					Username: "test-user",
					Email:    "test@user.com",
					Password: "P@ssw0rd123!",
				})
				assert.NoError(u, err)
				<-receiver

				user, err := userClient.Get("test-user")
				assert.NoError(u, err)

				user.Pending = false
				_, err = userClient.Status(user)
				assert.NoError(u, err)
			},
			Send: func(u *testing.T) {
				user, err := userClient.Get("test-user")
				assert.NoError(u, err)

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
		},
	} {
		t.Run(test.Name, func(u *testing.T) {
			test.Prepare(u)
			test.Send(u)

			e := <-receiver

			e.User.ResourceVersion = ""

			assert.NoError(u, err)
			assert.Equal(u, test.ExpectedType, e.Type)
			assert.Equal(u, test.ExpectedData, e.User)
		})
	}
}
