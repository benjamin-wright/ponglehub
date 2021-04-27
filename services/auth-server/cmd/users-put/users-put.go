package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/server"
)

type UserPut struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Verified *bool  `json:"verified"`
}

func main() {
	server.Run(func(cli *client.AuthClient, r *gin.Engine) {
		r.PUT("/:id", func(c *gin.Context) {
			var body UserPut
			if err := c.ShouldBindJSON(&body); err != nil {
				logrus.Errorf("Error reading user data: %+v", err)
				c.Status(400)
				return
			}

			logrus.Infof("Updating user %s: \"%s\" \"%s\" \"%s\"", c.Param("id"), body.Email, body.Name, body.Password)

			err := cli.UpdateUser(c.Request.Context(), c.Param("id"), client.User{
				Name:     body.Name,
				Email:    body.Email,
				Password: body.Password,
			}, body.Verified)

			if err != nil {
				logrus.Errorf("Error updating user: %+v", err)
				c.Status(500)
				return
			}

			c.Status(202)
		})
	})
}
