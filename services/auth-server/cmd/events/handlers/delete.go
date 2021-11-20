package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/lib/user-events/pkg/events"
)

func DeleteUser(db *client.PostgresClient, user events.User) {
	id := user.ID
	if _, err := uuid.Parse(id); err != nil {
		logrus.Warnf("Failed to delete user with badly formed id: %s", id)
		return
	}

	found, err := db.DeleteUser(context.TODO(), id)
	if err != nil {
		logrus.Errorf("Error deleting user \"%s\": %+v", id, err)
		return
	}

	if !found {
		logrus.Warnf("Failed to delete user \"%s\": Not found", id)
		return
	}
}
