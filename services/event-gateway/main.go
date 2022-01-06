package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/events/gateway/internal/managers/server"
	"ponglehub.co.uk/events/gateway/internal/managers/state"
	"ponglehub.co.uk/events/gateway/internal/services/crds"
	"ponglehub.co.uk/events/gateway/internal/services/tokens"
	"ponglehub.co.uk/events/gateway/internal/services/user_store"
)

func getEnv(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		logrus.Fatalf("Cannot find environment variable: %s", env)
	}

	return value
}

func getServices() (*crds.UserClient, *tokens.Tokens, *user_store.Store) {
	keyFilePath := getEnv("KEY_FILE")
	redisUrl := getEnv("REDIS_URL")

	crds.AddToScheme(scheme.Scheme)
	client, err := crds.New(&crds.ClientArgs{})
	if err != nil {
		logrus.Fatalf("Failed to start user client: %+v", err)
	}

	tokens, err := tokens.New(keyFilePath, redisUrl)
	if err != nil {
		logrus.Fatalf("Failed to start server: %+v", err)
	}

	store := user_store.New()

	return client, tokens, store
}

func main() {
	logrus.Infof("Starting operator...")

	client, tokens, store := getServices()

	stopServer := server.Start("BROKER_URL", getEnv("TOKEN_DOMAIN"), client, store, tokens)
	defer stopServer()

	stopListener := state.Start(client, store, tokens)
	defer stopListener()

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")
}
