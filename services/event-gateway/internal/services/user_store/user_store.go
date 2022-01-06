package user_store

import "github.com/sirupsen/logrus"

type Store struct {
	users map[string]string
}

func New() *Store {
	return &Store{
		users: map[string]string{},
	}
}

func (u *Store) Add(userId string, email string) {
	logrus.Infof("loading user %s", email)
	id, ok := u.users[email]
	if ok && id != userId {
		logrus.Errorf("user %s already exists in lookup!", email)
		return
	}

	u.users[email] = userId
}

func (u *Store) Remove(email string) {
	logrus.Infof("unloading user %s", email)

	delete(u.users, email)
}

func (u *Store) GetID(email string) (string, bool) {
	id, ok := u.users[email]

	return id, ok
}
