package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/operators/db/internal/client"
	"ponglehub.co.uk/operators/db/internal/types"
)

func main() {
	logrus.Infof("Starting operator...")

	client.AddToScheme(scheme.Scheme)

	cli, err := client.New()
	if err != nil {
		logrus.Fatalf("Failed to start operator client: %+v", err)
	}

	_, clientStopper := cli.ClientListen(
		func(newClient types.Client) {
			logrus.Infof("adding client: %+v", newClient)
		},
		func(oldClient types.Client, newClient types.Client) {
			logrus.Infof("updating client: %+v -> %+v", oldClient, newClient)
		},
		func(oldClient types.Client) {
			logrus.Infof("deleting client: %+v", oldClient)
		},
	)

	_, dbStopper := cli.DBListen(
		func(newDB types.Database) {
			logrus.Infof("adding database: %+v", newDB)
		},
		func(oldDB types.Database, newDB types.Database) {
			logrus.Infof("updating database: %+v -> %+v", oldDB, newDB)
		},
		func(oldDB types.Database) {
			logrus.Infof("deleteting database: %+v", oldDB)
		},
	)

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")

	clientStopper <- struct{}{}
	dbStopper <- struct{}{}

	log.Println("Stopped")
}
