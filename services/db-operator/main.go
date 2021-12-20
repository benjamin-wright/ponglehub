package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
	"ponglehub.co.uk/operators/db/internal/manager"
	"ponglehub.co.uk/operators/db/internal/types"
)

var addDeployment = manager.AddDeployment
var deleteDeployment = manager.DeleteDeployment
var addClient = manager.AddClient
var deleteClient = manager.DeleteClient

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

	_, clientStopper := crdClient.ClientListen(
		func(newClient types.Client) {
			logrus.Infof("Adding client: %s:%s (%s)", newClient.Database, newClient.Name, newClient.Namespace)

			err := addClient(crdClient, newClient)
			if err != nil {
				logrus.Errorf("Adding client failed: %+v", err)
				return
			}

			logrus.Infof("Client %s for database %s (%s) added", newClient.Name, newClient.Database, newClient.Namespace)
		},
		func(oldClient types.Client, newClient types.Client) {
			logrus.Infof("Updating client: %+v -> %+v", oldClient, newClient)
		},
		func(oldClient types.Client) {
			logrus.Infof("Deleting client: %s:%s (%s)", oldClient.Database, oldClient.Name, oldClient.Namespace)

			deleteClient(deplClient, oldClient)

			logrus.Infof("Client %s:%s (%s) deleted", oldClient.Database, oldClient.Name, oldClient.Namespace)
		},
	)

	_, dbStopper := crdClient.DBListen(
		func(newDB types.Database) {
			logrus.Infof("Adding database: %+v", newDB)

			err := addDeployment(deplClient, newDB)
			if err != nil {
				logrus.Errorf("Adding DB failed: %+v", err)
				return
			}

			logrus.Infof("Database %s (%s) added", newDB.Name, newDB.Namespace)
		},
		func(oldDB types.Database, newDB types.Database) {
			logrus.Infof("Updating database: %+v -> %+v", oldDB, newDB)
		},
		func(oldDB types.Database) {
			logrus.Infof("Deleteting database %s (%s)", oldDB.Name, oldDB.Namespace)

			deleteDeployment(deplClient, oldDB)

			logrus.Infof("Database %s (%s) deleted", oldDB.Name, oldDB.Namespace)
		},
	)

	_, deplStopper := deplClient.Listen(
		func(name string, namespace string, ready bool) {
			if err := crdClient.DBUpdate(name, namespace, ready); err != nil {
				logrus.Errorf("Failed to update CRD status: %s (%s) - %+v", name, namespace, err)
			}
		},
	)

	logrus.Infof("Running...")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server...")

	clientStopper <- struct{}{}
	dbStopper <- struct{}{}
	deplStopper <- struct{}{}

	log.Println("Stopped")
}
