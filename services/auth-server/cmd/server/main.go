package main

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
	"ponglehub.co.uk/auth/auth-server/cmd/server/routes"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

func main() {
	port, ok := os.LookupEnv("PONGLE_SERVER_PORT")
	if !ok {
		port = "80"
	}

	config, err := connect.ConfigFromEnv()
	if err != nil {
		logrus.Fatalf("Failed to load db config from environment: %+v", err)
	}

	cli, err := client.NewPostgresClient(context.Background(), config)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %+v", err)
	}

	r := gin.Default()

	r.POST("/", routes.LoginHandler(cli))

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
	quit := make(chan os.Signal, 1)
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
