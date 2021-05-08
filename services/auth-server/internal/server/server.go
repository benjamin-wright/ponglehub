package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/internal/client"
)

func GetRouter(database string, builder func(cli *client.AuthClient, r *gin.Engine)) *gin.Engine {
	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		logrus.Fatal("Environment Variable DB_HOST not found")
	}

	username, ok := os.LookupEnv("DB_USER")
	if !ok {
		username = "authserver"
	}

	cli, err := client.New(context.Background(), &client.AuthClientConfig{
		Username: username,
		Host:     host,
		Port:     26257,
		Database: database,
	})

	if err != nil {
		logrus.Fatalf("Failed to connect to database: %+v", err)
	}

	r := gin.Default()

	builder(cli, r)

	return r
}

func Run(builder func(cli *client.AuthClient, r *gin.Engine)) {
	port, ok := os.LookupEnv("PONGLE_SERVER_PORT")
	if !ok {
		logrus.Fatal("Environment Variable PONGLE_SERVER_PORT not found")
	}

	database, ok := os.LookupEnv("DB_NAME")
	if !ok {
		logrus.Fatal("Environment Variable DB_NAME not found")
	}

	r := GetRouter(database, builder)

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", port),
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
