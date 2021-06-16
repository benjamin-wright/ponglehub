package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/server"
)

type UserPost struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RouteBuilder(cli *client.AuthClient, r *gin.Engine) {
	r.POST("/", func(c *gin.Context) {
		var body UserPost
		if err := c.ShouldBindJSON(&body); err != nil {
			logrus.Errorf("Error reading user data: %+v", err)
			c.Status(400)
			return
		}

		logrus.Infof("Adding user: %s %s", body.Email, body.Name)

		hashed, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			logrus.Errorf("Error hashing user password: %+v", err)
			c.Status(500)
			return
		}

		userId, err := cli.AddUser(c.Request.Context(), client.User{
			Name:     body.Name,
			Email:    body.Email,
			Password: string(hashed),
		})

		if err != nil {
			logrus.Errorf("Error adding user: %+v", err)
			c.Status(500)
			return
		}

		if userId != "" {
			c.Status(202)
			c.JSON(202, gin.H{
				"id": userId,
			})
		} else {
			c.Status(409)
		}
	})
}

func main() {
	server.Run(RouteBuilder)
}
