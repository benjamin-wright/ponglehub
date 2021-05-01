package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/server"
)

func main() {
	server.Run(func(cli *client.AuthClient, r *gin.Engine) {
		r.GET("/:id", func(c *gin.Context) {
			user, err := cli.GetUser(c.Request.Context(), c.Param("id"))
			if err != nil {
				logrus.Errorf("Error getting user: %+v", err)
				c.Status(500)
				return
			}

			if user == nil {
				logrus.Warnf("User \"%s\" not found", c.Param("id"))
				c.Status(404)
				return
			}

			c.JSON(200, gin.H{
				"name":     user.Name,
				"email":    user.Email,
				"verified": user.Verified,
			})
		})
	})
}
