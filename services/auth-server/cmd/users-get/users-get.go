package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/server"
)

func RouteBuilder(cli server.AuthClient, r *gin.Engine) {
	r.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			logrus.Warnf("Failed to delete user with badly formed id: %s", id)
			c.Status(400)
			return
		}

		user, err := cli.GetUser(c.Request.Context(), id)
		if err != nil {
			logrus.Errorf("Error getting user: %+v", err)
			c.Status(500)
			return
		}

		if user == nil {
			logrus.Warnf("User \"%s\" nots found", id)
			c.Status(404)
			return
		}

		c.JSON(200, gin.H{
			"name":     user.Name,
			"email":    user.Email,
			"verified": user.Verified,
		})
	})
}

func main() {
	server.Run(RouteBuilder)
}
