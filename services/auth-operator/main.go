package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/auth/auth-operator/internal/handlers"
	"ponglehub.co.uk/auth/auth-operator/internal/users"
	events "ponglehub.co.uk/lib/user-events"
)

func main() {
	logrus.Infof("Starting operator...")

	users.AddToScheme(scheme.Scheme)
	userClient, err := users.New(&users.ClientArgs{})
	if err != nil {
		logrus.Fatalf("Failed to start user client: %+v", err)
	}

	eventClient, err := events.New("BROKER_URL", "auth-operator")
	if err != nil {
		logrus.Fatalf("Failed to start event client: %+v", err)
	}

	handler, err := handlers.New(eventClient, userClient)
	if err != nil {
		logrus.Fatalf("Failed to create handler: %+v", err)
	}

	_, stopper := userClient.Listen(handler.AddUser, handler.UpdateUser, handler.DeleteUser)
	defer func() { stopper <- struct{}{} }()

	cancelFunc, err := events.Listen(80, handler.UserEvent)
	if err != nil {
		logrus.Fatalf("Failed to create event listener: %+v", err)
	}
	defer cancelFunc()

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")
}
