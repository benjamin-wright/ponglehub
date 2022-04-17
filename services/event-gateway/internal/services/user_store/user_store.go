package user_store

import "github.com/sirupsen/logrus"

type Store struct {
	idLookup   map[string]string
	nameLookup map[string]string
}

func New() *Store {
	return &Store{
		idLookup:   map[string]string{},
		nameLookup: map[string]string{},
	}
}

func (u *Store) Add(userId string, name string, email string) {
	logrus.Infof("loading user %s", email)
	id, ok := u.idLookup[email]
	if ok && id != userId {
		logrus.Errorf("can't add user %s, already exists in lookup!", email)
		return
	}

	u.idLookup[email] = userId
	u.nameLookup[userId] = name
}

func (u *Store) Remove(email string) {
	logrus.Infof("unloading user %s", email)

	id, ok := u.idLookup[email]
	if !ok {
		logrus.Errorf("can't delete user %s, doesn't exist in lookup!", email)
		return
	}

	delete(u.nameLookup, id)
	delete(u.idLookup, email)
}

func (u *Store) GetID(email string) (string, bool) {
	id, ok := u.idLookup[email]

	return id, ok
}

func (u *Store) GetName(id string) (string, bool) {
	name, ok := u.nameLookup[id]

	return name, ok
}

func (u *Store) ListIDs(id string) []string {
	logrus.Infof("lookup: %+v", u.nameLookup)
	numIds := len(u.nameLookup) - 1
	if numIds < 1 {
		return []string{}
	}

	ids := make([]string, numIds)
	i := 0
	for x := range u.nameLookup {
		if id == x {
			continue
		}

		if i < numIds {
			ids[i] = x
			i++
		} else {
			ids = append(ids, x)
		}
	}

	return ids
}
