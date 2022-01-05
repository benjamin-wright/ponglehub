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
	users  map[string]string
}

func Start(brokerEnv string, domain string, tokens *tokens.Tokens, crds *crds.UserClient) (*Server, error) {
	server := Server{
		users:  map[string]string{},
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

	r.POST("/events", func(c *gin.Context) {
		_, err := c.Cookie("ponglehub.login")
		if err == http.ErrNoCookie {
			c.Status(401)
			return
		}

		if err != nil {
			logrus.Errorf("Error getting cookie: %+v", err)
			c.Status(500)
			return
		}

		h.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/auth/login", func(c *gin.Context) {
	})

	type LoginBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	r.POST("/auth/login", func(c *gin.Context) {
		body := LoginBody{}
		c.Bind(&body)

		id, ok := server.users[body.Email]
		if !ok {
			logrus.Errorf("Login user not found: %s", body.Email)
			c.Status(400)
			return
		}

		ok, err := tokens.CheckPassword(id, body.Password)
		if err != nil {
			logrus.Errorf("Failed checking user password: %+v", err)
			c.Status(500)
			return
		}

		if !ok {
			logrus.Errorf("Passwords didn't match for user %s", body.Email)
			c.Status(400)
			return
		}

		token, err := tokens.NewToken(id, "login", 1*time.Hour)
		if err != nil {
			logrus.Errorf("Failed creating token for user %s: %+v", body.Email, err)
			c.Status(500)
			return
		}

		c.SetCookie("ponglehub.login", token, 6400, "/", domain, false, true)
		c.Status(200)
	})

	r.GET("/auth/set-password", func(c *gin.Context) {
	})

	type SetPasswordBody struct {
		Invite   string `json:"invite"`
		Password string `json:"password"`
		Confirm  string `json:"confirm"`
	}

	r.POST("/auth/set-password", func(c *gin.Context) {
		body := SetPasswordBody{}
		c.Bind(&body)

		if body.Password != body.Confirm {
			logrus.Errorf("Mismatched password and confirmation")
			c.JSON(400, gin.H{"failure": "passwords"})
			return
		}

		claims, err := tokens.Parse(body.Invite)
		if err != nil {
			logrus.Errorf("Failed to parse invite token: %+v", err)
			c.JSON(400, gin.H{"failure": "token"})
			return
		}

		if claims.Kind != "invite" {
			logrus.Errorf("Tried to set password without an invite token: %s", claims.Kind)
			c.Status(401)
			return
		}

		err = tokens.AddPasswordHash(claims.Subject, body.Password)
		if err != nil {
			logrus.Errorf("Failed to hash password: %+v", err)
			c.Status(500)
			return
		}

		logrus.Infof("Password updated for user %s", claims.Subject)
	})

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

func (s *Server) addToLookup(user crds.User) {
	logrus.Infof("loading user %s", user.Email)
	id, ok := s.users[user.Email]
	if ok && id != user.ID {
		logrus.Errorf("user %s already exists in lookup!", user.Email)
		return
	}

	s.users[user.Email] = user.ID
}

func (s *Server) removeFromLookup(user crds.User) {
	logrus.Infof("unloading user %s", user.Email)

	delete(s.users, user.Email)
}

func (s *Server) AddUser(newUser crds.User) {
	newUser = s.processUser(newUser)
	s.addToLookup(newUser)
}

func (s *Server) UpdateUser(oldUser crds.User, newUser crds.User) {
	newUser = s.processUser(newUser)

	if oldUser.Email != newUser.Email {
		s.removeFromLookup(oldUser)
		s.addToLookup(newUser)
	}
}

func (s *Server) RemoveUser(oldUser crds.User) {
	s.removeFromLookup(oldUser)
	if oldUser.Invited {
		err := s.tokens.DeleteToken(oldUser.ID, "invite")
		if err != nil {
			logrus.Errorf("Error removing user: %+v")
		}
	}
}
