package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/games/draughts/cmd/server/server"
	"ponglehub.co.uk/games/draughts/pkg/database"
	"ponglehub.co.uk/lib/events"
)

func main() {
	logrus.Infof("Starting operator...")

	client, err := events.New(events.EventsArgs{
		BrokerEnv: "BROKER_URL",
		Source:    "naughts-and-crosses",
	})
	if err != nil {
		logrus.Fatalf("failed to create events client: %+v", err)
	}

	db, err := database.New()
	if err != nil {
		logrus.Fatalf("failed to create database client: %+v", err)
	}

	server, err := server.New(client, db)
	if err != nil {
		logrus.Fatalf("Failed to start server: %+v", err)
	}
	defer server.Stop()

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logrus.Infof("Stopped")
}
