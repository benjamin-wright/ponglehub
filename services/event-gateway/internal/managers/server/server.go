package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/events/gateway/internal/services/tokens"
	"ponglehub.co.uk/events/gateway/internal/services/user_store"
	"ponglehub.co.uk/events/gateway/pkg/crds"
	"ponglehub.co.uk/lib/events"
)

func Start(brokerEnv string, domain string, origins []string, crdClient *crds.UserClient, store *user_store.Store, tokens *tokens.Tokens) func() {
	eventClient, err := events.New(events.EventsArgs{BrokerEnv: brokerEnv, Source: "event-gateway"})
	if err != nil {
		logrus.Fatalf("Failed to create broker client: %+v", err)
	}

	engine := gin.Default()

	engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://ponglehub.co.uk",
			"http://games.ponglehub.co.uk",
		},
		AllowMethods: []string{"POST", "GET"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "ce-dataschema",
			"ce-id", "ce-source", "ce-specversion",
			"ce-subject", "ce-time", "ce-type",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	engine.LoadHTMLGlob("/html/*")

	engine.GET("/events", eventsGetRoute(tokens, domain, store, eventClient))
	engine.GET("/auth/user", userRoute(tokens, domain, crdClient, store))
	engine.GET("/auth/login", loginHTML)
	engine.POST("/auth/login", loginRoute(store, tokens, domain))
	engine.POST("/auth/logout", logoutRoute(tokens, domain))
	engine.GET("/auth/set-password", setPasswordHTML)
	engine.POST("/auth/set-password", setPasswordRoute(store, crdClient, tokens))

	server := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: engine,
	}

	go func() {
		// service connections
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatalf("Error starting external server: %+v\n", err)
		}
	}()

	return func() {
		err := server.Close()
		if err != nil {
			logrus.Errorf("Error closing server: %+v", err)
		}
	}
}

func watchEvents(conn *websocket.Conn) (<-chan cloudevents.Event, <-chan struct{}) {
	events := make(chan cloudevents.Event)
	stopper := make(chan struct{})

	go func(events chan<- cloudevents.Event, stopper chan<- struct{}) {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				logrus.Warnf("Closing websocket connection: %+v", err)
				stopper <- struct{}{}
				return
			}

			var event cloudevents.Event
			err = json.Unmarshal(msg, &event)
			if err != nil {
				logrus.Errorf("Error unmarshalling event data: %+v", err)
				continue
			}

			events <- event
		}
	}(events, stopper)

	return events, stopper
}

func eventsGetRoute(tokens *tokens.Tokens, domain string, store *user_store.Store, client *events.Events) func(c *gin.Context) {
	var wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	return func(c *gin.Context) {
		subject, err := loggedIn(c, tokens, domain)
		if err != nil {
			return
		}

		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			logrus.Errorf("Failed to set websocket upgrade: %+v", err)
			return
		}

		responses, stopper, err := tokens.WatchResponses(subject)
		if err != nil {
			logrus.Errorf("Failed to watch responses: %+v", err)
			return
		}

		events, stopped := watchEvents(conn)

		for {
			select {
			case <-stopped:
				stopper <- struct{}{}
				return
			case response := <-responses:
				logrus.Infof("sending response to %s", subject)
				err = conn.WriteMessage(websocket.TextMessage, []byte(response))
				if err != nil {
					logrus.Errorf("Error return response: %+v", err)
				}
			case event := <-events:
				switch event.Type() {
				case "auth.list-friends":
					logrus.Infof("listing friends for: %s", subject)
					friends := store.GetFriends(subject)

					response, err := json.Marshal(map[string]interface{}{
						"type": "auth.list-friends.response",
						"data": friends,
					})
					if err != nil {
						logrus.Errorf("Error serialising list-friends response: %+v", err)
						continue
					}

					err = conn.WriteMessage(websocket.TextMessage, response)
					if err != nil {
						logrus.Errorf("Error returning list-friends response: %+v", err)
					}

				default:
					logrus.Infof("passing through event: %s", event.Type())

					event.SetExtension("userid", subject)
					err = client.Proxy(event)
					if err != nil {
						logrus.Errorf("Error proxying event to broker: %+v", err)
					}
				}
			}
		}
	}
}

func loggedIn(c *gin.Context, tokens *tokens.Tokens, domain string) (string, error) {
	token, err := c.Cookie("ponglehub.login")
	if err == http.ErrNoCookie {
		c.Status(http.StatusUnauthorized)
		return "", err
	}

	if err != nil {
		logrus.Errorf("Error getting cookie: %+v", err)
		c.Status(http.StatusInternalServerError)
		return "", err
	}

	claims, err := tokens.Parse(token)
	if err != nil {
		logrus.Errorf("Error parsing cookie: %+v", err)
		c.Status(http.StatusUnauthorized)
		return "", err
	}

	if claims.Kind != "login" {
		logrus.Errorf("Accessed with non login cookie: %s", claims.Kind)
		c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
		c.Status(http.StatusUnauthorized)
		return "", errors.New("something")
	}

	return claims.Subject, nil
}

func userRoute(tokens *tokens.Tokens, domain string, users *crds.UserClient, store *user_store.Store) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := c.Cookie("ponglehub.login")
		if err == http.ErrNoCookie {
			logrus.Errorf("No cookie found")
			c.Status(http.StatusUnauthorized)
			return
		}

		if err != nil {
			logrus.Errorf("Error getting cookie: %+v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		claims, err := tokens.Parse(token)
		if err != nil {
			logrus.Errorf("Error parsing cookie: %+v", err)
			c.Status(http.StatusUnauthorized)
			return
		}

		if claims.Kind != "login" {
			logrus.Errorf("Accessed with non login cookie: %s", claims.Kind)
			c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
			c.Status(http.StatusUnauthorized)
			return
		}

		t, err := tokens.GetToken(claims.Subject, "login")
		if err != nil {
			logrus.Errorf("Failed to fetch invite token: %+v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if t == "" {
			logrus.Errorf("Invite token expired: %s", claims.Subject)
			c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
			c.Status(http.StatusUnauthorized)
			return
		}

		if t != token {
			logrus.Errorf("Login token doesn't match latest: %s", claims.Subject)
			c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
			c.Status(http.StatusUnauthorized)
			return
		}

		name, ok := store.GetName(claims.Subject)
		if !ok {
			logrus.Errorf("Failed to find user name: user not found")
			c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
			c.Status(http.StatusUnauthorized)
			return
		}

		user, err := users.Get(name)
		if err != nil {
			logrus.Errorf("Failed to fetch user data: %+v", err)
			c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
			c.Status(http.StatusUnauthorized)
			return
		}

		c.JSON(
			http.StatusOK,
			gin.H{
				"name": user.Display,
			},
		)
	}
}

func loginHTML(c *gin.Context) {
	url, ok := c.GetQuery("redirect")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"redirect": url,
	})
}

type loginBody struct {
	Email    string `json:"email" form:"email" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Redirect string `json:"redirect" form:"redirect" binding:"required"`
}

func loginRoute(store *user_store.Store, tokens *tokens.Tokens, domain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		body := loginBody{}
		c.Bind(&body)

		if body.Email == "" || body.Password == "" || body.Redirect == "" {
			logrus.Errorf("Missing login params")
			c.JSON(http.StatusBadRequest, gin.H{"failure": "bad input"})
			return
		}

		id, ok := store.GetID(body.Email)
		if !ok {
			logrus.Errorf("Login user not found: %s", body.Email)

			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"redirect": body.Redirect,
				"error":    true,
			})
			return
		}

		ok, err := tokens.CheckPassword(id, body.Password)
		if err != nil {
			logrus.Errorf("Failed checking user password: %+v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if !ok {
			logrus.Errorf("Wrong password for user %s", body.Email)

			c.HTML(http.StatusOK, "login.tmpl", gin.H{
				"redirect": body.Redirect,
				"error":    true,
			})
			return
		}

		token, err := tokens.NewToken(id, "login", 1*time.Hour)
		if err != nil {
			logrus.Errorf("Failed creating token for user %s: %+v", body.Email, err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.SetCookie("ponglehub.login", token, 6400, "/", domain, false, true)
		c.Redirect(http.StatusFound, body.Redirect)
	}
}

func logoutRoute(tokens *tokens.Tokens, domain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, err := c.Cookie("ponglehub.login")
		if err == http.ErrNoCookie {
			logrus.Errorf("No cookie found")
			c.Status(http.StatusUnauthorized)
			return
		}

		if err != nil {
			logrus.Errorf("Error getting cookie: %+v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		claims, err := tokens.Parse(token)
		if err != nil {
			logrus.Errorf("Error parsing cookie: %+v", err)
			c.Status(http.StatusUnauthorized)
			return
		}

		if claims.Kind != "login" {
			logrus.Errorf("Accessed with non login cookie: %s", claims.Kind)
			c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
			c.Status(http.StatusUnauthorized)
			return
		}

		t, err := tokens.GetToken(claims.Subject, "login")
		if err != nil {
			logrus.Errorf("Failed to fetch invite token: %+v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if t == "" {
			logrus.Errorf("Invite token expired: %s", claims.Subject)
			c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
			c.Status(http.StatusUnauthorized)
			return
		}

		if t != token {
			logrus.Errorf("Login token doesn't match latest: %s", claims.Subject)
			c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
			c.Status(http.StatusUnauthorized)
			return
		}

		err = tokens.DeleteToken(claims.Subject, "login")
		if err != nil {
			logrus.Errorf("Failed revoking token for user %s: %+v", claims.Subject, err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.SetCookie("ponglehub.login", "", 0, "/", domain, false, true)
		c.Status(http.StatusNoContent)
	}
}

func setPasswordHTML(c *gin.Context) {
	token, ok := c.GetQuery("token")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	c.HTML(http.StatusOK, "set-password.tmpl", gin.H{
		"invite": token,
	})
}

type setPasswordBody struct {
	Invite   string `json:"invite" form:"invite" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Confirm  string `json:"confirm" form:"confirm" binding:"required"`
}

func setPasswordRoute(store *user_store.Store, crdClient *crds.UserClient, tokens *tokens.Tokens) func(c *gin.Context) {
	return func(c *gin.Context) {
		body := setPasswordBody{}
		c.Bind(&body)

		if body.Invite == "" || body.Password == "" || body.Confirm == "" {
			logrus.Errorf("Missing set password params")
			c.JSON(http.StatusBadRequest, gin.H{"failure": "bad input"})
			return
		}

		if body.Password != body.Confirm {
			logrus.Errorf("Mismatched password and confirmation")
			c.JSON(http.StatusBadRequest, gin.H{"failure": "passwords"})
			return
		}

		claims, err := tokens.Parse(body.Invite)
		if err != nil {
			logrus.Errorf("Failed to parse invite token: %+v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"failure": "token"})
			return
		}

		if claims.Kind != "invite" {
			logrus.Errorf("Tried to set password without an invite token: %s", claims.Kind)
			c.Status(http.StatusUnauthorized)
			return
		}

		t, err := tokens.GetToken(claims.Subject, "invite")
		if err != nil {
			logrus.Errorf("Failed to fetch invite token: %+v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if t == "" {
			logrus.Errorf("Invite token expired: %s", claims.Subject)
			c.Status(http.StatusUnauthorized)
			return
		}

		if t != body.Invite {
			logrus.Errorf("Invite token doesn't match latest: %s", claims.Subject)
			c.Status(http.StatusUnauthorized)
			return
		}

		err = tokens.AddPasswordHash(claims.Subject, body.Password)
		if err != nil {
			logrus.Errorf("Failed to hash password: %+v", err)
			c.Status(http.StatusInternalServerError)
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
