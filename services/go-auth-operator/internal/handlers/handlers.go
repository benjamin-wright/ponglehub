package handlers

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-operator/internal/users"
)

type Handler struct {
	events Events
	users  Users
}

type Events interface {
	NewUser(user users.User) error
	UpdateUser(user users.User) error
	DeleteUser(user users.User) error
}

type Users interface {
	Update(user users.User) error
	Status(user users.User) error
	Get(name string) (users.User, error)
}

func New(eventClient Events, userClient Users) (*Handler, error) {
	return &Handler{
		events: eventClient,
		users:  userClient,
	}, nil
}

func (h *Handler) AddUser(user users.User) {
	if user.Pending {
		logrus.Infof("Not adding user '%s', already pending", user.Name)
		return
	}

	if user.ID != "" {
		logrus.Infof("Not adding user '%s', already got ID '%s", user.Name, user.ID)
		return
	}

	logrus.Infof("Setting status to pending for %s", user.Name)
	user.Pending = true
	err := h.users.Status(user)
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", user.Name, err)
		return
	}

	logrus.Infof("Sending add user event '%s'", user.Name)
	err = h.events.NewUser(user)
	if err == nil {
		return
	}
}

func (h *Handler) UpdateUser(oldUser users.User, newUser users.User) {
	if oldUser.Equals(newUser) {
		logrus.Infof("Not updating using '%s': No spec changes", newUser.Name)
		return
	}

	if newUser.Pending {
		logrus.Infof("Not updating user '%s': Already pending", newUser.Name)
		return
	}

	newUser.Pending = true
	err := h.users.Status(newUser)
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", newUser.Name, err)
		return
	}

	logrus.Infof("Sending update user event '%s'", newUser.Name)
	err = h.events.UpdateUser(newUser)
	if err == nil {
		logrus.Errorf("Failed to update user %s: %+v", newUser.Name, err)
		return
	}
}

func (h *Handler) DeleteUser(user users.User) {
	logrus.Infof("Sending delete user event '%s'", user.Name)
	err := h.events.DeleteUser(user)
	if err != nil {
		logrus.Errorf("Failed to delete user %s: %+v", user.Name, err)
	}
}

func (h *Handler) SetUser(name string, id string) {
	user, err := h.users.Get(name)
	if err != nil {
		logrus.Errorf("Failed to fetch user %s: %+v", name, err)
		return
	}

	user.ID = id
	user.Pending = false
	err = h.users.Status(user)
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", name, err)
		return
	}
}
