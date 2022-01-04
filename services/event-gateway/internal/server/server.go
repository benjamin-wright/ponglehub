package server

import (
	"context"
	"fmt"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/events/gateway/internal/crds"
	"ponglehub.co.uk/events/gateway/internal/tokens"
	"ponglehub.co.uk/lib/events"
)

type Server struct {
	cancel context.CancelFunc
	client *events.Events
	tokens *tokens.Tokens
	crds   *crds.UserClient
	users  map[string]user
}

func Start(brokerEnv string, tokens *tokens.Tokens, crds *crds.UserClient) (*Server, error) {
	server := Server{
		users:  map[string]user{},
		crds:   crds,
		tokens: tokens,
	}

	client, err := events.New(events.EventsArgs{
		BrokerEnv: brokerEnv,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create broker client: %+v", err)
	}

	cancelFunc, err := events.Listen(80, server.handle)
	if err != nil {
		return nil, fmt.Errorf("failed to create event server: %+v", err)
	}

	server.client = client
	server.cancel = cancelFunc

	return &server, nil
}

func (s *Server) Stop() {
	s.cancel()
}

func (s *Server) handle(ctx context.Context, event event.Event) {
	tokenString, err := s.tokens.NewToken("id-1234")
	if err != nil {
		logrus.Errorf("Error proxying event to broker: %+v", err)
		return
	} else {
		logrus.Infof("Made a token: %s", tokenString)
	}

	err = s.client.Proxy(event)
	if err != nil {
		logrus.Errorf("Error proxying event to broker: %+v", err)
	}
}

type user struct {
	id           string
	passwordHash string
}

func (s *Server) conditionalInvite(user crds.User) crds.User {
	if user.PasswordHash != "" || user.InviteToken != "" {
		return user
	}

	logrus.Infof("issuing invite token for %s", user.Email)

	tokenString, err := s.tokens.NewToken("id-1234")
	if err != nil {
		logrus.Errorf("Error proxying event to broker: %+v", err)
		return user
	}

	user.InviteToken = tokenString
	_, err = s.crds.Update(user)
	if err != nil {
		logrus.Errorf("Error updating invite token: %+v", err)
		user.InviteToken = ""
	}

	return user
}

func (s *Server) AddUser(newUser crds.User) {
	if _, ok := s.users[newUser.Email]; ok {
		logrus.Infof("reloading user %s", newUser.Email)
	} else {
		logrus.Infof("loading user %s", newUser.Email)
	}

	newUser = s.conditionalInvite(newUser)

	s.users[newUser.Email] = user{
		id:           newUser.ID,
		passwordHash: newUser.PasswordHash,
	}
}

func (s *Server) UpdateUser(oldUser crds.User, newUser crds.User) {
	s.AddUser(newUser)
}

func (s *Server) RemoveUser(oldUser crds.User) {
	logrus.Infof("unloading user %s", oldUser.Email)

	delete(s.users, oldUser.Email)
}
