package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"

	"ponglehub.co.uk/events/responder/internal/storage"
	"ponglehub.co.uk/lib/events"
)

func getEnv(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		logrus.Fatalf("Cannot find environment variable: %s", env)
	}

	return value
}

func main() {
	redisUrl := getEnv("REDIS_URL")
	store, err := storage.New(redisUrl)
	if err != nil {
		logrus.Fatalf("Failed to create storage client: %+v", err)
	}

	cancelFunc, err := events.Listen(80, func(ctx context.Context, event event.Event) {
		userIdObj, err := event.Context.GetExtension("userid")
		if err != nil {
			logrus.Errorf("Failed to get user id from event: %+v", err)
			return
		}

		userId, ok := userIdObj.(string)
		if !ok {
			logrus.Errorf("Expected user id to be a string, got %T", userId)
			return
		}

		if userId == "" {
			logrus.Infof("Not responding to event %s, empty userId", event.Type())
			return
		}

		err = store.AddEvent(userId, event)
		if err != nil {
			logrus.Errorf("Failed to store event %s: %+v", event.Type(), err)
			return
		}

		logrus.Infof("Stored event: %s", event.Type())
	})
	if err != nil {
		logrus.Fatalf("Failed to start event listener: %+v", err)
	}
	defer cancelFunc()

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")
}
