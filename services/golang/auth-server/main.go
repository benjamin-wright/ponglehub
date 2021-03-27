package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth-server/internal/client"
)

type UserPost struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	r := gin.Default()
	cli, err := client.New(context.Background(), &client.AuthClientConfig{
		Username: "authserver",
		Host:     "auth-server-cockroach-public",
		Port:     26257,
		Database: "authserver",
	})

	if err != nil {
		logrus.Fatalf("Failed to connect to database: %+v", err)
	}

	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	r.GET("/user", func(c *gin.Context) {
		users, err := cli.ListUsers(c.Request.Context())
		if err != nil {
			logrus.Errorf("Error getting list of users: %+v", err)
			c.Status(500)
			return
		}

		c.JSON(200, users)
	})

	r.GET("/user/:user", func(c *gin.Context) {
		user, err := cli.GetUser(c.Request.Context(), c.Param("user"))
		if err != nil {
			logrus.Errorf("Error getting user: %+v", err)
			c.Status(500)
			return
		}

		if user == nil {
			logrus.Warnf("User \"%s\" not found", c.Param("user"))
			c.Status(404)
			return
		}

		c.JSON(200, gin.H{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"verified": user.Verified,
		})
	})

	r.POST("/user", func(c *gin.Context) {
		var body UserPost
		if err := c.ShouldBindJSON(&body); err != nil {
			logrus.Errorf("Error reading user data user: %+v", err)
			c.Status(400)
			return
		}

		logrus.Infof("Adding user: %s %s %s", body.Email, body.Name, body.Password)

		c.Status(202)
	})

	r.Run("0.0.0.0:80")
}
