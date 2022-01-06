package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/events/gateway/internal/crds"
	"ponglehub.co.uk/events/gateway/internal/tokens"
	"ponglehub.co.uk/lib/events"
)

type Server struct {
	cancel func()
	client *events.Events
	tokens *tokens.Tokens
	crds   *crds.UserClient
	store  *UserStore
}

func Start(brokerEnv string, domain string, tokens *tokens.Tokens, crds *crds.UserClient) (*Server, error) {
	server := Server{
		store:  NewUserStore(),
		crds:   crds,
		tokens: tokens,
	}

	client, err := events.New(events.EventsArgs{
		BrokerEnv: brokerEnv,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create broker client: %+v", err)
	}

	ctx := context.Background()
	p, err := cloudevents.NewHTTP()
	if err != nil {
		logrus.Fatalf("failed to create protocol: %s", err.Error())
	}

	h, err := cloudevents.NewHTTPReceiveHandler(ctx, p, server.handle)
	if err != nil {
		logrus.Fatalf("failed to create handler: %s", err.Error())
	}

	r := gin.Default()

	r.POST("/events", EventsRoute(h))
	r.GET("/auth/login", func(c *gin.Context) {})
	r.POST("/auth/login", LoginRoute(server.store, tokens, domain))
	r.GET("/auth/set-password", func(c *gin.Context) {})
	r.POST("/auth/set-password", SetPasswordRoute(tokens))

	srv := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatalf("Error starting server: %+v\n", err)
		}
	}()

	server.client = client
	server.cancel = func() {
		err := srv.Close()
		if err != nil {
			logrus.Errorf("Error closing server: %+v", err)
		}
	}

	return &server, nil
}

func (s *Server) Stop() {
	s.cancel()
}

func (s *Server) handle(ctx context.Context, event event.Event) {
	logrus.Infof("passing through event: %s", event.Type())
	err := s.client.Proxy(event)
	if err != nil {
		logrus.Errorf("Error proxying event to broker: %+v", err)
	}
}

func (s *Server) setUserStatus(user crds.User) {
	_, err := s.crds.Status(user)
	if err != nil {
		logrus.Errorf("Error updating user status: %+v", err)
	}
}

func (s *Server) processUser(user crds.User) crds.User {
	password, err := s.tokens.GetToken(user.ID, "password")
	if err != nil {
		logrus.Errorf("Error fetching password: %+v", err)
		return user
	}

	if password != "" {
		if user.Invited || !user.Member {
			logrus.Infof("restoring status for member %s", user.Email)
			user.Invited = false
			user.Member = true
			s.setUserStatus(user)
		}

		return user
	}

	invite, err := s.tokens.GetToken(user.ID, "invite")
	if err != nil {
		logrus.Errorf("Error fetching invite token: %+v", err)
		return user
	}

	if invite != "" {
		if !user.Invited || user.Member {
			logrus.Infof("restoring status for invited user %s", user.Email)
			user.Invited = true
			user.Member = false
			s.setUserStatus(user)
		}

		return user
	}

	logrus.Infof("issuing invite token for %s", user.Email)

	_, err = s.tokens.NewToken(user.ID, "invite", 72*time.Hour)
	if err != nil {
		logrus.Errorf("Error creating invite token: %+v", err)
		return user
	}

	if !user.Invited || user.Member {
		logrus.Infof("setting status for invited user %s", user.Email)
		user.Invited = true
		user.Member = false
		s.setUserStatus(user)
	}

	return user
}

func (s *Server) AddUser(newUser crds.User) {
	newUser = s.processUser(newUser)
	s.store.Add(newUser.ID, newUser.Email)
}

func (s *Server) UpdateUser(oldUser crds.User, newUser crds.User) {
	newUser = s.processUser(newUser)

	if oldUser.Email != newUser.Email {
		s.store.Remove(oldUser.Email)
		s.store.Add(newUser.ID, newUser.Email)
	}
}

func (s *Server) RemoveUser(oldUser crds.User) {
	s.store.Remove(oldUser.Email)
	if oldUser.Invited {
		err := s.tokens.DeleteToken(oldUser.ID, "invite")
		if err != nil {
			logrus.Errorf("Error removing user: %+v")
		}
	}
}
