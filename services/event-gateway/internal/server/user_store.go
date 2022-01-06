package server

import "github.com/sirupsen/logrus"

type UserStore struct {
	users map[string]string
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: map[string]string{},
	}
}

func (u *UserStore) Add(userId string, email string) {
	logrus.Infof("loading user %s", email)
	id, ok := u.users[email]
	if ok && id != userId {
		logrus.Errorf("user %s already exists in lookup!", email)
		return
	}

	u.users[email] = userId
}

func (u *UserStore) Remove(email string) {
	logrus.Infof("unloading user %s", email)

	delete(u.users, email)
}

func (u *UserStore) GetID(email string) (string, bool) {
	id, ok := u.users[email]

	return id, ok
}
