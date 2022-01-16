package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/games/naughts-and-crosses/cmd/server/server"
)

func main() {
	logrus.Infof("Starting operator...")

	stopper, err := server.Start()
	if err != nil {
		logrus.Fatalf("Failed to start server: %+v", err)
	}

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")

	stopper()

	log.Println("Stopped")
}
