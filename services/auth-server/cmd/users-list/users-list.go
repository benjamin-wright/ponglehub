package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/server"
)

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

func RouteBuilder(cli client.AuthClient, r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		users, err := cli.ListUsers(c.Request.Context())
		if err != nil {
			logrus.Errorf("Error getting list of users: %+v", err)
			c.Status(500)
			return
		}

		responses := []UserResponse{}
		for _, user := range users {
			responses = append(responses, UserResponse{
				ID:       user.ID,
				Name:     user.Name,
				Email:    user.Email,
				Verified: user.Verified,
			})
		}

		c.JSON(200, responses)
	})
}

func main() {
	server.Run(RouteBuilder)
}
