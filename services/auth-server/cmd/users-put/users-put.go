package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/server"
)

type UserPut struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Verified *bool  `json:"verified"`
}

func (u UserPut) HasValues() bool {
	return u.Name != "" || u.Email != "" || u.Password != "" || u.Verified != nil
}

func RouteBuilder(cli *client.AuthClient, r *gin.Engine) {
	r.PUT("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			logrus.Warnf("Failed to delete user with badly formed id: %s", id)
			c.Status(400)
			return
		}

		var body UserPut
		if err := c.ShouldBindJSON(&body); err != nil {
			logrus.Warnf("Error reading user data for id %s: %+v", id, err)
			c.Status(400)
			return
		}

		if !body.HasValues() {
			logrus.Warnf("No update values provided for %s", id)
			c.Status(400)
			return
		}

		logrus.Infof("Updating user %s: \"%s\" \"%s\" \"%t\"", id, body.Email, body.Name, body.Password != "")

		hashedPassword := ""
		if body.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
			if err != nil {
				logrus.Errorf("Error hashing user password: %+v", err)
				c.Status(500)
				return
			}

			hashedPassword = string(hash)
		}

		success, err := cli.UpdateUser(c.Request.Context(), id, client.User{
			Name:     body.Name,
			Email:    body.Email,
			Password: hashedPassword,
		}, body.Verified)

		if err != nil {
			logrus.Errorf("Error updating user: %+v", err)
			c.Status(500)
			return
		}

		if !success {
			logrus.Errorf("Error updating user: %s not found", id)
			c.Status(404)
			return
		}

		c.Status(202)
	})
}

func main() {
	server.Run(RouteBuilder)
}
