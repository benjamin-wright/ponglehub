package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/auth/auth-server/cmd/events/handlers"
	"ponglehub.co.uk/auth/auth-server/internal/client"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
	"ponglehub.co.uk/lib/user-events/pkg/events"
)

func main() {
	brokerUrl, ok := os.LookupEnv("BROKER_URL")
	if !ok {
		logrus.Fatal("Environment Variable BROKER_URL not found")
	}

	cfg, err := connect.ConfigFromEnv()
	if err != nil {
		logrus.Fatalf("Failed to load db config from environment: %+v", err)
	}

	db, err := client.NewPostgresClient(context.Background(), cfg)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %+v", err)
	}

	client, err := events.New(brokerUrl)
	if err != nil {
		logrus.Fatalf("Failed to create event client: %+v", err)
	}

	cancel, err := events.Listen(func(event events.UserEvent) {
		switch event.Type {
		case "ponglehub.auth.user.delete":
			logrus.Infof("Got user delete event for %s", event.User.Name)
			handlers.DeleteUser(db, event.User)
		case "ponglehub.auth.user.add":
			logrus.Infof("Got user add event for %s", event.User.Name)
			handlers.AddUser(db, client, event.User)
		case "ponglehub.auth.user.update":
			logrus.Infof("Got user update event for %s", event.User.Name)
			handlers.UpdateUser(db, client, event.User)
		default:
			logrus.Warnf("Unrecognised event type: %s", event.Type)
		}
	})
	if err != nil {
		logrus.Fatalf("Failed to start event listener: %+v", err)
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")
	cancel()
	log.Println("Server exiting")
}
