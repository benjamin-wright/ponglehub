package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"ponglehub.co.uk/auth/auth-server/internal/server"
)

type LoginPost struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func routeBuilder(cli server.AuthClient, r *gin.Engine) {
	r.POST("/", func(c *gin.Context) {
		var body LoginPost
		if err := c.ShouldBindJSON(&body); err != nil {
			logrus.Errorf("Error reading login data: %+v", err)
			c.Status(400)
			return
		}

		logrus.Infof("Fetching user: %s", body.Email)

		user, err := cli.GetUserByEmail(c.Request.Context(), body.Email)
		if err != nil {
			logrus.Errorf("Error getting user: %+v", err)
			c.Status(500)
			return
		}

		if user == nil {
			logrus.Infof("User %s failed login: user not found", body.Email)
			c.Status(401)
			return
		}

		if !user.Verified {
			logrus.Infof("User %s failed login: user not yet verified", body.Email)
			c.Status(401)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err != nil {
			logrus.Warnf("User failed password check: %+v", err)
			c.Status(401)
			return
		}

		c.JSON(200, gin.H{
			"id":   user.ID,
			"name": user.Name,
		})
	})
}

func main() {
	server.Run(routeBuilder)
}
