package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/lib/user-events/pkg/events"
)

func UpdateUser(db *client.PostgresClient, sender *events.Events, user events.User) {
	id := user.ID
	if _, err := uuid.Parse(id); err != nil {
		logrus.Warnf("Failed to delete user with badly formed id: %s", id)
		return
	}

	logrus.Infof("Updating user %s: \"%s\" \"%s\" \"%t\"", id, user.Email, user.Username, user.Password != "")

	hashedPassword := ""
	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			logrus.Errorf("Error hashing user password: %+v", err)
			return
		}

		hashedPassword = string(hash)
	}

	verified := true

	success, err := db.UpdateUser(context.TODO(), id, client.User{
		Name:     user.Username,
		Email:    user.Email,
		Password: hashedPassword,
	}, &verified)

	if err != nil {
		logrus.Errorf("Error updating user: %+v", err)
		return
	}

	if !success {
		logrus.Errorf("Error updating user: %s not found", id)
	}
}
