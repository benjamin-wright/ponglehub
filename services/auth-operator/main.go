package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/auth/auth-operator/internal/handlers"
	"ponglehub.co.uk/auth/auth-operator/internal/users"
	"ponglehub.co.uk/lib/user-events/pkg/events"
)

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		logrus.Fatalf("Environment variable value required: %s", key)
	}

	return value
}

func main() {
	logrus.Infof("Starting operator...")

	brokerUrl := getEnv("BROKER_URL")

	users.AddToScheme(scheme.Scheme)

	userClient, err := users.New()
	if err != nil {
		logrus.Fatalf("Failed to start user client: %+v", err)
	}

	eventClient, err := events.New(brokerUrl)
	if err != nil {
		logrus.Fatalf("Failed to start event client: %+v", err)
	}

	handler, err := handlers.New(eventClient, userClient)
	if err != nil {
		logrus.Fatalf("Failed to create handler: %+v", err)
	}

	_, stopper := userClient.Listen(handler.AddUser, handler.UpdateUser, handler.DeleteUser)
	cancelEventListener, err := events.Listen(handler.UserEvent)
	if err != nil {
		logrus.Fatalf("Failed to create event listener: %+v", err)
	}

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")

	stopper <- struct{}{}
	cancelEventListener()

	log.Println("Stop signal sent")
}
