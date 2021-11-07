package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/auth/auth-operator/internal/client"
	"ponglehub.co.uk/auth/auth-operator/internal/handlers"
)

func main() {
	logrus.Infof("Starting operator...")

	client.AddToScheme(scheme.Scheme)

	cli, err := client.New()
	if err != nil {
		logrus.Fatalf("Failed to start client: %+v", err)
	}

	handler, err := handlers.New()
	if err != nil {
		logrus.Fatalf("Failed to create handler: %+v", err)
	}

	_, stopper := cli.Listen(handler.AddUser, handler.UpdateUser, handler.DeleteUser)

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")

	stopper <- struct{}{}

	log.Println("Stop signal sent")
}
