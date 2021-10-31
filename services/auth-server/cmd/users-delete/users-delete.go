package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/server"
)

func RouteBuilder(cli server.AuthClient, r *gin.Engine) {
	r.DELETE("/:id", func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			logrus.Warnf("Failed to delete user with badly formed id: %s", id)
			c.Status(400)
			return
		}

		found, err := cli.DeleteUser(c.Request.Context(), id)
		if err != nil {
			logrus.Errorf("Error deleting user \"%s\": %+v", id, err)
			c.Status(500)
			return
		}

		if !found {
			logrus.Warnf("Failed to delete user \"%s\": Not found", id)
			c.Status(404)
			return
		}

		c.Status(204)
	})
}

func main() {
	server.Run(RouteBuilder)
}
