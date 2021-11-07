package handlers

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-operator/internal/client"
	"ponglehub.co.uk/auth/auth-operator/internal/events"
)

type Handler struct {
	events *events.Events
}

func New() (*Handler, error) {
	client, err := events.New()
	if err != nil {
		return nil, err
	}

	return &Handler{
		events: client,
	}, nil
}

func (h *Handler) AddUser(user *client.AuthUser) {
	logrus.Infof("Adding user %s", user.ObjectMeta.Name)
	err := h.events.NewUser(user)
	if err != nil {
		logrus.Errorf("Failed to add user %s: %+v", user.ObjectMeta.Name, err)
	}
}

func (h *Handler) UpdateUser(oldUser *client.AuthUser, newUser *client.AuthUser) {
	logrus.Infof("Updating user %s", newUser.ObjectMeta.Name)
	err := h.events.UpdateUser(newUser)
	if err != nil {
		logrus.Errorf("Failed to update user %s: %+v", newUser.ObjectMeta.Name, err)
	}
}

func (h *Handler) DeleteUser(user *client.AuthUser) {
	logrus.Infof("Deleting user %s", user.ObjectMeta.Name)
	err := h.events.DeleteUser(user)
	if err != nil {
		logrus.Errorf("Failed to delete user %s: %+v", user.ObjectMeta.Name, err)
	}
}
