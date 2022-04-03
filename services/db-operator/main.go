package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/database"
	"ponglehub.co.uk/operators/db/internal/deployments"
	"ponglehub.co.uk/operators/db/internal/reconcilers"
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

	dbClient := database.New()

	events := make(chan interface{}, 5)

	databaseStore, dbStopper := crdClient.DBListen(events)
	clientStore, clientStopper := crdClient.ClientListen(events)
	statefulSetStore, statefulsetStopper := deplClient.ListenStatefulSets(events)
	serviceStore, serviceStopper := deplClient.ListenService(events)

	dbReconciler := reconcilers.NewDeploymentReconciler(
		crdClient,
		deplClient,
		databaseStore,
		statefulSetStore,
		serviceStore,
	)
	dbReconcilerStopper := dbReconciler.Start(events)

	clientReconciler := reconcilers.NewClientReconciler(
		crdClient,
		deplClient,
		dbClient,
		databaseStore,
		clientStore,
		statefulSetStore,
	)
	clientReconcilerStopper := clientReconciler.Start(events)

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
	dbReconcilerStopper <- struct{}{}
	clientReconcilerStopper <- struct{}{}

	log.Println("Stopped")
}
