package main

import (
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/games/naughts-and-crosses/cmd/server/routes"
	"ponglehub.co.uk/games/naughts-and-crosses/pkg/database"
	"ponglehub.co.uk/lib/events"
)

func main() {
	logrus.Infof("Starting operator...")

	db, err := database.New()
	if err != nil {
		logrus.Fatalf("failed to create database client: %+v", err)
	}

	stop, err := events.Serve(events.ServeParams{
		BrokerEnv: "BROKER_URL",
		Source:    "naughts-and-crosses",
		Routes: events.EventRoutes{
			"naughts-and-crosses.list-games": routes.ListGames(db),
			"naughts-and-crosses.new-game":   routes.NewGame(db),
			"naughts-and-crosses.load-game":  routes.LoadGame(db),
			"naughts-and-crosses.mark":       routes.Mark(db),
		},
	})
	if err != nil {
		logrus.Fatalf("failed to start server: %+v", err)
	}
	defer stop()

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logrus.Infof("Stopped")
}
