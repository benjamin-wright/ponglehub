package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/lib/events"
)

func getString(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		logrus.Fatalf("Environment variable %s not defined", env)
	}

	return value
}

func getInt(env string) int {
	value := getString(env)

	intValue, err := strconv.Atoi(value)
	if err != nil {
		logrus.Fatalf("Failed to load environment variable %s, error converting to int: %+v", env, err)
	}

	return intValue
}

func main() {
	eventPort := getInt("EVENT_PORT")
	serverPort := getInt("SERVER_PORT")

	eventList := []event.Event{}

	cancelFunc, err := events.Listen(eventPort, func(ctx context.Context, event event.Event) {
		logrus.Infof("Recording event: %s", event.Type())
		eventList = append(eventList, event)
	})
	if err != nil {
		logrus.Fatalf("Failed to start event listener: %+v", err)
	}

	r := gin.Default()

	r.POST("/clear", func(c *gin.Context) {
		eventList = []event.Event{}
		c.Status(200)
	})

	r.GET("/events", func(c *gin.Context) {
		types := []string{}
		for _, event := range eventList {
			types = append(types, event.Type())
		}

		c.JSON(200, types)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", serverPort),
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			logrus.Fatalf("Error starting server: %+v\n", err)
		}
	}()

	logrus.Infof("Running server: %d, and events: %d...", serverPort, eventPort)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")

	cancelFunc()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")

	log.Println("Stopped")
}
