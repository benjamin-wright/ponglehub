package main

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/games/draughts/pkg/database"
	"ponglehub.co.uk/games/draughts/pkg/routes"
	"ponglehub.co.uk/lib/events"
)

func main() {
	logrus.Infof("Starting server...")

	db, err := database.New()
	if err != nil {
		logrus.Fatalf("failed to create database client: %+v", err)
	}

	events.Serve(events.ServeParams{
		BrokerEnv: "BROKER_URL",
		Source:    "draughts",
		Routes: events.EventRoutes{
			"draughts.list-games": routes.ListGames(db),
			"draughts.new-game":   routes.NewGame(db),
			"draughts.load-game":  nil,
		},
	})
}
