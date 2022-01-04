package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/events/broker/internal/crds"
	"ponglehub.co.uk/events/broker/internal/router"
	"ponglehub.co.uk/events/broker/internal/server"
)

func main() {
	logrus.Infof("Starting operator...")

	r := router.New()

	crds.AddToScheme(scheme.Scheme)
	crdClient, err := crds.New(&crds.ClientArgs{})
	if err != nil {
		logrus.Fatalf("Failed to start operator client: %+v", err)
	}

	_, crdStopper := crdClient.Listen(func(oldTrigger *crds.EventTrigger, newTrigger *crds.EventTrigger) {
		logrus.Infof("Detected trigger change")

		if oldTrigger != nil {
			for _, filter := range oldTrigger.Spec.Filters {
				if err := r.Remove(filter, oldTrigger.Spec.URL); err != nil {
					logrus.Errorf("failed to remove %s -> %s: %+v", filter, oldTrigger.Spec.URL, err)
				} else {
					logrus.Infof("removed %s -> %s", filter, oldTrigger.Spec.URL)
				}
			}
		}

		if newTrigger != nil {
			for _, filter := range newTrigger.Spec.Filters {
				r.Add(filter, newTrigger.Spec.URL)
				logrus.Infof("added %s -> %s", filter, newTrigger.Spec.URL)
			}
		}
	})

	serverStopper, err := server.Start(&r)
	if err != nil {
		logrus.Fatalf("Failed to listen for events: %+v", err)
	}

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")

	crdStopper <- struct{}{}
	serverStopper()

	log.Println("Stopped")
}
