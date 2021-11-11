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

func getEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		logrus.Fatalf("Environment variable value required: %s", key)
	}

	return value
}

func main() {
	logrus.Infof("Starting operator...")

	natsUrl := getEnv("NATS_URL")
	natsSubject := getEnv("NATS_SUBJECT")

	client.AddToScheme(scheme.Scheme)

	cli, err := client.New()
	if err != nil {
		logrus.Fatalf("Failed to start client: %+v", err)
	}

	handler, err := handlers.New(natsUrl, natsSubject)
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
