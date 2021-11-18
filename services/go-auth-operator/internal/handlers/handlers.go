package handlers

import (
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ponglehub.co.uk/auth/auth-operator/internal/client"
)

type Handler struct {
	events Events
	users  Users
}

type Events interface {
	NewUser(user *client.AuthUser) error
	UpdateUser(user *client.AuthUser) error
	DeleteUser(user *client.AuthUser) error
}

type Users interface {
	SetStatus(user *client.AuthUser, opts v1.UpdateOptions) error
	Get(name string, opts v1.GetOptions) (*client.AuthUser, error)
}

func New(eventClient Events, userClient Users) (*Handler, error) {
	return &Handler{
		events: eventClient,
		users:  userClient,
	}, nil
}

func (h *Handler) AddUser(user *client.AuthUser) {
	if user == nil {
		logrus.Error("Nil passed to AddUser handler")
		return
	}

	if user.Status.Pending {
		logrus.Infof("Not adding user '%s', already pending", user.Name)
		return
	}

	if user.Status.ID != "" {
		logrus.Infof("Not adding user '%s', already got ID '%s", user.Name, user.Status.ID)
		return
	}

	logrus.Infof("Setting status to pending for %s", user.Name)
	user.Status.Pending = true
	err := h.users.SetStatus(user, v1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", user.Name, err)
		return
	}

	logrus.Infof("Sending add user event '%s'", user.ObjectMeta.Name)
	err = h.events.NewUser(user)
	if err == nil {
		return
	}

	logrus.Errorf("Failed to add user %s: %+v", user.Name, err)

	latest, err := h.users.Get(user.Name, v1.GetOptions{})
	if err != nil {
		logrus.Errorf("Failed to fetch latest user %s: %+v", user.Name, err)
		return
	}

	latest.Status.Pending = false
	err = h.users.SetStatus(latest, v1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", user.Name, err)
		return
	}
}

func (h *Handler) UpdateUser(oldUser *client.AuthUser, newUser *client.AuthUser) {
	if oldUser == nil || newUser == nil {
		logrus.Errorf("Nil passed to UpdateUser handler")
		return
	}

	if oldUser.Equals(newUser) {
		logrus.Infof("Not updating using '%s': No spec changes", newUser.Name)
		return
	}

	if newUser.Status.Pending {
		logrus.Infof("Not updating user '%s': Already pending", newUser.Name)
		return
	}

	newUser.Status.Pending = true
	err := h.users.SetStatus(newUser, v1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", newUser.Name, err)
		return
	}

	logrus.Infof("Sending update user event '%s'", newUser.Name)
	err = h.events.UpdateUser(newUser)
	if err == nil {
		return
	}

	logrus.Errorf("Failed to update user %s: %+v", newUser.Name, err)

	latest, err := h.users.Get(newUser.Name, v1.GetOptions{})
	if err != nil {
		logrus.Errorf("Failed to fetch latest user %s: %+v", newUser.Name, err)
		return
	}

	latest.Status.Pending = false
	err = h.users.SetStatus(latest, v1.UpdateOptions{})
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", newUser.Name, err)
		return
	}
}

func (h *Handler) DeleteUser(user *client.AuthUser) {
	if user == nil {
		logrus.Error("Nil passed to DeleteUser handler")
		return
	}

	logrus.Infof("Sending delete user event '%s'", user.Name)
	err := h.events.DeleteUser(user)
	if err != nil {
		logrus.Errorf("Failed to delete user %s: %+v", user.Name, err)
	}
}
