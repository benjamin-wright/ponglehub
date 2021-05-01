package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/auth/auth-server/internal/server"
)

func main() {
	server.Run(func(cli *client.AuthClient, r *gin.Engine) {
		r.DELETE("/:id", func(c *gin.Context) {
			found, err := cli.DeleteUser(c.Request.Context(), c.Param("id"))
			if err != nil {
				logrus.Errorf("Error deleting user \"%s\": %+v", c.Param("id"), err)
				c.Status(500)
				return
			}

			if !found {
				logrus.Warnf("Failed to delete user \"%s\": Not found", c.Param("id"))
				c.Status(404)
				return
			}

			c.Status(204)
		})
	})
}
