package main

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/games/naughts-and-crosses/pkg/database"
	"ponglehub.co.uk/games/naughts-and-crosses/pkg/routes"
	"ponglehub.co.uk/lib/events"
)

func main() {
	logrus.Infof("Starting server...")

	db, err := database.New()
	if err != nil {
		logrus.Fatalf("failed to create database client: %+v", err)
	}

	err = events.Serve(events.ServeParams{
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
		logrus.Fatal(err)
	}
}
