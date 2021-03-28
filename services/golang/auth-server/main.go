package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth-server/internal/client"
)

type UserPost struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
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

	r.GET("/user/:name", func(c *gin.Context) {
		user, err := cli.GetUser(c.Request.Context(), c.Param("name"))
		if err != nil {
			logrus.Errorf("Error getting user: %+v", err)
			c.Status(500)
			return
		}

		if user == nil {
			logrus.Warnf("User \"%s\" not found", c.Param("name"))
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
			logrus.Errorf("Error reading user data: %+v", err)
			c.Status(400)
			return
		}

		logrus.Infof("Adding user: %s %s %s", body.Email, body.Name, body.Password)

		err = cli.AddUser(c.Request.Context(), client.User{
			Name:     body.Name,
			Email:    body.Email,
			Password: body.Password,
		})
		if err != nil {
			logrus.Errorf("Error adding user: %+v", err)
			c.Status(500)
			return
		}

		c.Status(202)
	})

	r.DELETE("/user/:name", func(c *gin.Context) {
		found, err := cli.DeleteUser(c.Request.Context(), c.Param("name"))
		if err != nil {
			logrus.Errorf("Error deleting user \"%s\"", c.Param("name"))
			c.Status(500)
			return
		}

		if !found {
			logrus.Warnf("Failed to delete user \"%s\": Not found", c.Param("name"))
			c.Status(404)
			return
		}

		c.Status(202)
	})

	srv := &http.Server{
		Addr:    "0.0.0.0:80",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
