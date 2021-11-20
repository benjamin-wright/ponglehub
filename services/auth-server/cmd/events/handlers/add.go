package handlers

import (
	"context"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/lib/user-events/pkg/events"
)

func AddUser(db *client.PostgresClient, sender *events.Events, user events.User) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("Error hashing user password: %+v", err)
		return
	}

	userId, err := db.AddUser(context.TODO(), client.User{
		Name:     user.Username,
		Email:    user.Email,
		Password: string(hashed),
		Verified: false,
	})

	if err != nil {
		logrus.Errorf("Error adding user: %+v", err)
		return
	}

	user.ID = userId
	if err = sender.SetUser(user); err != nil {
		logrus.Errorf("Error sending set-user event: %+v", err)
	}
}
