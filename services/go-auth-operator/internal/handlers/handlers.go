package handlers

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-operator/internal/events"
)

type Handler struct {
	events Events
	users  Users
}

type Events interface {
	NewUser(user events.User) error
	UpdateUser(user events.User) error
	DeleteUser(user events.User) error
}

type Users interface {
	Update(user events.User) (events.User, error)
	Status(user events.User) (events.User, error)
	Get(name string) (events.User, error)
}

func New(eventClient Events, userClient Users) (*Handler, error) {
	return &Handler{
		events: eventClient,
		users:  userClient,
	}, nil
}

func (h *Handler) AddUser(user events.User) {
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
	_, err := h.users.Status(user)
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", user.Name, err)
		return
	}

	logrus.Infof("Sending add user event '%s'", user.Name)
	err = h.events.NewUser(user)
	if err != nil {
		logrus.Errorf("Failed to send event %s: %+v", user.Name, err)
		return
	}
}

func (h *Handler) UpdateUser(oldUser events.User, newUser events.User) {
	if oldUser.Equals(newUser) {
		logrus.Infof("Not updating user '%s': No spec changes", newUser.Name)
		return
	}

	if newUser.Pending {
		logrus.Infof("Not updating user '%s': Already pending", newUser.Name)
		return
	}

	newUser.Pending = true
	logrus.Infof("Setting status to pending for %s", newUser.Name)
	_, err := h.users.Status(newUser)
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", newUser.Name, err)
		return
	}

	logrus.Infof("Sending update user event '%s'", newUser.Name)
	err = h.events.UpdateUser(newUser)
	if err != nil {
		logrus.Errorf("Failed to update user %s: %+v", newUser.Name, err)
		return
	}
}

func (h *Handler) DeleteUser(user events.User) {
	logrus.Infof("Sending delete user event '%s'", user.Name)
	err := h.events.DeleteUser(user)
	if err != nil {
		logrus.Errorf("Failed to delete user %s: %+v", user.Name, err)
	}
}

func (h *Handler) UserEvent(event events.UserEvent) {
	if event.Type != "ponglehub.auth.user.set" {
		logrus.Warnf("Unrecognised event type: %s", event)
		return
	}

	user := event.User

	currentUser, err := h.users.Get(user.Name)
	if err != nil {
		logrus.Errorf("Failed to fetch user %s: %+v", user.Name, err)
		return
	}

	if currentUser.ResourceVersion != user.ResourceVersion {
		logrus.Infof("Ignoring user set event: resource version %s -> %s", user.ResourceVersion, currentUser.ResourceVersion)
		return
	}

	user.Pending = false
	_, err = h.users.Status(user)
	if err != nil {
		logrus.Errorf("Failed to update user status %s: %+v", user.Name, err)
		return
	}
}
