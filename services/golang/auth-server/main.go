package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth-server/internal/client"
)

func main() {
	r := gin.Default()
	cli, err := client.New(context.Background(), &client.AuthClientConfig{
		Username: "user",
		Host:     "host",
		Port:     1234,
	})

	if err != nil {
		logrus.Fatalf("Failed to connect to database: %+v", err)
	}

	r.GET("/status", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	r.GET("/users", func(c *gin.Context) {
		users, err := cli.ListUsers(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Aw snap",
			})
		}

		c.JSON(200, users)
	})

	r.Run("0.0.0.0:80")
}
