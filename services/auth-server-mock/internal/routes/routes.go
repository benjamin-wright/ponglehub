package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	userState "ponglehub.co.uk/auth-server-mock/internal/state"
)

func RouteBuilder(r *gin.Engine, state *userState.State) {
	r.GET("/users", func(c *gin.Context) {
		responses := []gin.H{}

		for _, user := range state.Users {
			responses = append(responses, gin.H{
				"id":       user.ID,
				"name":     user.Name,
				"email":    user.Email,
				"verified": user.Verified,
			})
		}

		c.JSON(200, responses)
	})

	r.POST("/users", func(c *gin.Context) {
		var body struct {
			Name     string
			Email    string
			Password string
		}

		if err := c.ShouldBindJSON(&body); err != nil {
			logrus.Errorf("Error reading user data: %+v", err)
			c.Status(400)
			return
		}

		logrus.Infof("Adding user: %s %s", body.Email, body.Name)

		ID := uuid.New()

		state.Users = append(state.Users, userState.User{
			ID:       ID.String(),
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
			Verified: false,
		})

		c.JSON(202, gin.H{
			"id": ID,
		})
	})
}
