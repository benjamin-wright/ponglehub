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
	"ponglehub.co.uk/lib/postgres/pkg/connect"
)

type AuthClient interface {
	AddUser(ctx context.Context, user client.User) (string, error)
	UpdateUser(ctx context.Context, id string, user client.User, verified *bool) (bool, error)
	GetUser(ctx context.Context, id string) (*client.User, error)
	GetUserByEmail(ctx context.Context, email string) (*client.User, error)
	ListUsers(ctx context.Context) ([]*client.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
}

func GetRouter(config connect.ConnectConfig, builder func(cli AuthClient, r *gin.Engine)) *gin.Engine {
	cli, err := client.NewPostgresClient(context.Background(), config)

	if err != nil {
		logrus.Fatalf("Failed to connect to database: %+v", err)
	}

	r := gin.Default()

	builder(cli, r)

	return r
}

func Run(builder func(cli AuthClient, r *gin.Engine)) {
	port, ok := os.LookupEnv("PONGLE_SERVER_PORT")
	if !ok {
		logrus.Fatal("Environment Variable PONGLE_SERVER_PORT not found")
	}

	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		logrus.Fatal("Environment Variable DB_HOST not found")
	}

	username, ok := os.LookupEnv("DB_USER")
	if !ok {
		username = "authserver"
	}

	password, ok := os.LookupEnv("DB_PASS")
	if !ok {
		logrus.Fatal("Enrivonment Variable DB_PASS not found")
	}

	certsDir, ok := os.LookupEnv("DB_CERTS")
	if !ok {
		logrus.Fatal("Enrivonment Variable DB_CERTS not found")
	}

	database, ok := os.LookupEnv("DB_NAME")
	if !ok {
		logrus.Fatal("Environment Variable DB_NAME not found")
	}

	r := GetRouter(connect.ConnectConfig{
		Username: username,
		Password: password,
		Host:     host,
		Port:     26257,
		Database: database,
		CertsDir: certsDir,
	}, builder)

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
