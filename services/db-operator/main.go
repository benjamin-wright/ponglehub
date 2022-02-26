package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
	"ponglehub.co.uk/operators/db/internal/reconciler"
)

func main() {
	logrus.Infof("Starting operator...")

	crds.AddToScheme(scheme.Scheme)

	crdClient, err := crds.New()
	if err != nil {
		logrus.Fatalf("Failed to start operator client: %+v", err)
	}

	deplClient, err := deployments.New()
	if err != nil {
		logrus.Fatalf("Failed to start operator client: %+v", err)
	}

	events := make(chan interface{}, 5)
	actions := make(chan interface{}, 5)

	_, clientStopper := crdClient.ClientListen(events)
	_, dbStopper := crdClient.DBListen(events)
	_, statefulsetStopper := deplClient.ListenStatefulSets(events)
	_, serviceStopper := deplClient.ListenService(events)
	reconcilerStopper := reconciler.Start(events, actions)

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")

	clientStopper <- struct{}{}
	dbStopper <- struct{}{}
	statefulsetStopper <- struct{}{}
	serviceStopper <- struct{}{}
	reconcilerStopper <- struct{}{}

	log.Println("Stopped")
}
