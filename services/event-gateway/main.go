package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/events/gateway/internal/crds"
	"ponglehub.co.uk/events/gateway/internal/server"
	"ponglehub.co.uk/events/gateway/internal/tokens"
)

func getEnv(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		logrus.Fatalf("Cannot find environment variable: %s", env)
	}

	return value
}

func main() {
	logrus.Infof("Starting operator...")

	keyFilePath := getEnv("KEY_FILE")

	crds.AddToScheme(scheme.Scheme)
	client, err := crds.New(&crds.ClientArgs{})
	if err != nil {
		logrus.Fatalf("Failed to start user client: %+v", err)
	}

	tk, err := tokens.New(keyFilePath)
	if err != nil {
		logrus.Fatalf("Failed to start server: %+v", err)
	}

	srv, err := server.Start("BROKER_URL", tk, client)
	if err != nil {
		logrus.Fatalf("Failed to start server: %+v", err)
	}
	defer srv.Stop()

	_, stopper := client.Listen(srv.AddUser, srv.UpdateUser, srv.RemoveUser)
	defer func() { stopper <- struct{}{} }()

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")
}
