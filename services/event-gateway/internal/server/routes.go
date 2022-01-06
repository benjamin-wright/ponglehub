package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/events/gateway/internal/tokens"
)

func EventsRoute(handler http.Handler) func(c *gin.Context) {
	return func(c *gin.Context) {
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

		handler.ServeHTTP(c.Writer, c.Request)
	}
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginRoute(store *UserStore, tokens *tokens.Tokens, domain string) func(c *gin.Context) {
	return func(c *gin.Context) {
		body := LoginBody{}
		c.Bind(&body)

		id, ok := store.GetID(body.Email)
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
	}
}

type SetPasswordBody struct {
	Invite   string `json:"invite"`
	Password string `json:"password"`
	Confirm  string `json:"confirm"`
}

func SetPasswordRoute(tokens *tokens.Tokens) func(c *gin.Context) {
	return func(c *gin.Context) {
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
	}
}
