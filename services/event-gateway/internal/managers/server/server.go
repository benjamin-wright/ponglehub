package server

import (
	"context"
	"net/http"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/events/gateway/internal/services/tokens"
	"ponglehub.co.uk/events/gateway/internal/services/user_store"
	"ponglehub.co.uk/events/gateway/pkg/crds"
	"ponglehub.co.uk/lib/events"
)

func Start(brokerEnv string, domain string, crdClient *crds.UserClient, store *user_store.Store, tokens *tokens.Tokens) func() {
	eventClient, err := events.New(events.EventsArgs{BrokerEnv: brokerEnv})
	if err != nil {
		logrus.Fatalf("Failed to create broker client: %+v", err)
	}

	r := gin.Default()

	r.LoadHTMLGlob("/html/*")

	r.POST("/events", eventsRoute(tokens, domain, eventClient))
	r.GET("/auth/login", func(c *gin.Context) {
		url, ok := c.GetQuery("redirect")
		if !ok {
			c.Status(400)
			return
		}

		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title":    "Main website",
			"redirect": url,
		})
	})
	r.POST("/auth/login", loginRoute(store, tokens, domain))
	r.GET("/auth/set-password", func(c *gin.Context) {})
	r.POST("/auth/set-password", setPasswordRoute(store, crdClient, tokens))

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

	return func() {
		err := srv.Close()
		if err != nil {
			logrus.Errorf("Error closing server: %+v", err)
		}
	}
}

func eventsRoute(tokens *tokens.Tokens, domain string, client *events.Events) func(c *gin.Context) {
	ctx := context.Background()
	p, err := cloudevents.NewHTTP()
	if err != nil {
		logrus.Fatalf("failed to create protocol: %s", err.Error())
	}

	handler := func(ctx context.Context, event event.Event) {
		logrus.Infof("passing through event: %s", event.Type())
		err := client.Proxy(event)
		if err != nil {
			logrus.Errorf("Error proxying event to broker: %+v", err)
		}
	}

	h, err := cloudevents.NewHTTPReceiveHandler(ctx, p, handler)
	if err != nil {
		logrus.Fatalf("failed to create handler: %s", err.Error())
	}

	return func(c *gin.Context) {
		token, err := c.Cookie("ponglehub.login")
		if err == http.ErrNoCookie {
			c.Status(401)
			return
		}

		if err != nil {
			logrus.Errorf("Error getting cookie: %+v", err)
			c.Status(500)
			return
		}

		claims, err := tokens.Parse(token)
		if err != nil {
			logrus.Errorf("Error parsing cookie: %+v", err)
			c.Status(401)
			return
		}

		if claims.Kind != "login" {
			logrus.Errorf("Accessed with non login cookie: %s", claims.Kind)
			c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
			c.Status(401)
			return
		}

		c.Request.Header.Add("ce-userid", claims.Subject)

		h.ServeHTTP(c.Writer, c.Request)
	}
}

type loginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginRoute(store *user_store.Store, tokens *tokens.Tokens, domain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		body := loginBody{}
		c.Bind(&body)

		if body.Email == "" || body.Password == "" {
			logrus.Errorf("Missing login params")
			c.JSON(400, gin.H{"failure": "bad input"})
			return
		}

		id, ok := store.GetID(body.Email)
		if !ok {
			logrus.Errorf("Login user not found: %s", body.Email)
			c.Status(401)
			return
		}

		ok, err := tokens.CheckPassword(id, body.Password)
		if err != nil {
			logrus.Errorf("Failed checking user password: %+v", err)
			c.Status(500)
			return
		}

		if !ok {
			logrus.Errorf("Wrong password for user %s", body.Email)
			c.Status(401)
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
	}
}

type setPasswordBody struct {
	Invite   string `json:"invite"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func setPasswordRoute(store *user_store.Store, crdClient *crds.UserClient, tokens *tokens.Tokens) func(c *gin.Context) {
	return func(c *gin.Context) {
		body := setPasswordBody{}
		c.Bind(&body)

		if body.Invite == "" || body.Password == "" || body.Confirm == "" {
			logrus.Errorf("Missing set password params")
			c.JSON(400, gin.H{"failure": "bad input"})
			return
		}

		if body.Password != body.Confirm {
			logrus.Errorf("Mismatched password and confirmation")
			c.JSON(400, gin.H{"failure": "passwords"})
			return
		}

		claims, err := tokens.Parse(body.Invite)
		if err != nil {
			logrus.Errorf("Failed to parse invite token: %+v", err)
			c.JSON(401, gin.H{"failure": "token"})
			return
		}

		if claims.Kind != "invite" {
			logrus.Errorf("Tried to set password without an invite token: %s", claims.Kind)
			c.Status(401)
			return
		}

		t, err := tokens.GetToken(claims.Subject, "invite")
		if err != nil {
			logrus.Errorf("Failed to fetch invite token: %+v", err)
			c.Status(500)
			return
		}

		if t == "" {
			logrus.Errorf("Invite token expired: %s", claims.Subject)
			c.Status(401)
			return
		}

		if t != body.Invite {
			logrus.Errorf("Invite token doesn't match latest: %s", claims.Subject)
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

		err = tokens.DeleteToken(claims.Subject, "invite")
		if err != nil {
			logrus.Errorf("Failed to delete token after setting password: %+v", err)
			return
		}

		name, ok := store.GetName(claims.Subject)
		if !ok {
			logrus.Errorf("Failed to update user after setting password: user not found")
			return
		}

		user, err := crdClient.Get(name)
		if err != nil {
			logrus.Errorf("Failed to update user after setting password: %+v", err)
			return
		}

		if user.Invited || !user.Member {
			user.Invited = false
			user.Member = true

			logrus.Infof("Updating status for new member %s", user.Email)
			_, err = crdClient.Status(user)
			if err != nil {
				logrus.Errorf("Failed to update user after setting password: %+v", err)
				return
			}
		}
	}
}
