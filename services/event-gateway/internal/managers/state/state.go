package state

import (
	"time"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/events/gateway/internal/services/tokens"
	"ponglehub.co.uk/events/gateway/internal/services/user_store"
	"ponglehub.co.uk/events/gateway/pkg/crds"
)

func Start(client *crds.UserClient, store *user_store.Store, tokens *tokens.Tokens) func() {
	_, stopper := client.Listen(
		func(newUser crds.User) {
			newUser = processUser(client, tokens, newUser)
			store.Add(newUser.ID, newUser.Name, newUser.Email)
		},
		func(oldUser crds.User, newUser crds.User) {
			newUser = processUser(client, tokens, newUser)

			if oldUser.Email != newUser.Email {
				store.Remove(oldUser.Email)
				store.Add(newUser.ID, newUser.Name, newUser.Email)
			}
		},
		func(oldUser crds.User) {
			store.Remove(oldUser.Email)
			if oldUser.Invited {
				err := tokens.DeleteToken(oldUser.ID, "invite")
				if err != nil {
					logrus.Errorf("Error removing user: %+v", err)
				}
			}
		},
	)

	return func() { stopper <- struct{}{} }
}

func setUserStatus(client *crds.UserClient, user crds.User) {
	_, err := client.Status(user)
	if err != nil {
		logrus.Errorf("Error updating user status: %+v", err)
	}
}

func processUser(client *crds.UserClient, tokens *tokens.Tokens, user crds.User) crds.User {
	password, err := tokens.GetToken(user.ID, "password")
	if err != nil {
		logrus.Errorf("Error fetching password: %+v", err)
		return user
	}

	if password != "" {
		if user.Invited || !user.Member {
			logrus.Infof("restoring status for member %s", user.Email)
			user.Invited = false
			user.Member = true
			setUserStatus(client, user)
		}

		return user
	}

	invite, err := tokens.GetToken(user.ID, "invite")
	if err != nil {
		logrus.Errorf("Error fetching invite token: %+v", err)
		return user
	}

	if invite != "" {
		if !user.Invited || user.Member {
			logrus.Infof("restoring status for invited user %s", user.Email)
			user.Invited = true
			user.Member = false
			setUserStatus(client, user)
		}

		return user
	}

	logrus.Infof("issuing invite token for %s", user.Email)

	_, err = tokens.NewToken(user.ID, "invite", 72*time.Hour)
	if err != nil {
		logrus.Errorf("Error creating invite token: %+v", err)
		return user
	}

	if !user.Invited || user.Member {
		logrus.Infof("setting status for invited user %s", user.Email)
		user.Invited = true
		user.Member = false
		setUserStatus(client, user)
	}

	return user
}
