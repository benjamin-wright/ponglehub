package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes/scheme"
	"ponglehub.co.uk/operators/db/internal/certs"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
	"ponglehub.co.uk/operators/db/internal/types"
)

func deleteDeployment(deplClient *deployments.DeploymentsClient, db types.Database) {
	if err := deplClient.DeleteDeployment(db); err != nil {
		logrus.Errorf("failed to delete database: %+v", err)
	}

	if err := deplClient.DeleteService(db.Namespace, db.Name); err != nil {
		logrus.Errorf("failed to delete service: %+v", err)
	}

	if err := deplClient.DeleteNodeSecret(db.Namespace, db.Name); err != nil {
		logrus.Errorf("failed to delete node secret: %+v", err)
	}

	if err := deplClient.DeleteCASecret(db.Namespace, db.Name); err != nil {
		logrus.Errorf("failed to delete ca secret: %+v", err)
	}
}

func addDeployment(deplClient *deployments.DeploymentsClient, db types.Database) error {
	if ok, err := deplClient.HasService(db.Namespace, db.Name); err != nil {
		return fmt.Errorf("failed getting service: %+v", err)
	} else if !ok {
		logrus.Infof("Creating service...")
		if err = deplClient.AddService(db.Namespace, db.Name); err != nil {
			return fmt.Errorf("failed adding service: %+v", err)
		}
	} else {
		logrus.Infof("Service already exists")
	}

	key, err := deplClient.GetCASecret(db.Namespace, db.Name)
	if err != nil {
		return fmt.Errorf("failed getting CA secret: %+v", err)
	}

	_, err = deplClient.GetNodeSecret(db.Namespace, db.Name)
	if err != nil {
		return fmt.Errorf("failed getting node secret: %+v", err)
	}

	if key == nil {
		logrus.Infof("Creating ssl secrets...")
		ca, caKey, err := certs.GenerateCACerts()
		if err != nil {
			return fmt.Errorf("failed to create ca cert: %+v", err)
		}

		err = deplClient.AddCaSecret(db.Namespace, db.Name, caKey)
		if err != nil {
			return fmt.Errorf("failed to create ca key secret: %+v", err)
		}

		dnsNames := []string{
			db.Name,
			fmt.Sprintf("%s.%s", db.Name, db.Namespace),
			fmt.Sprintf("%s.%s.svc.cluster.local", db.Name, db.Namespace),
		}

		node, nodeKey, err := certs.GenerateNodeCerts(dnsNames, ca, caKey)
		if err != nil {
			return fmt.Errorf("failed to create node cert: %+v", err)
		}

		err = deplClient.AddNodeSecret(db.Namespace, db.Name, deployments.NodeCerts{
			CACrt:   ca,
			NodeCrt: node,
			NodeKey: nodeKey,
		})
		if err != nil {
			return fmt.Errorf("failed to create node certs secret: %+v", err)
		}

		logrus.Infof("Created secrets")
	} else {
		logrus.Infof("Secrets already exist")
	}

	depls, err := deplClient.GetDeployments(db.Namespace)
	if err != nil {
		return fmt.Errorf("list DBs failed: %+v", err)
	}

	for _, depl := range depls {
		if depl.Name == db.Name {
			logrus.Infof("DB %s (%s) already exists", db.Name, db.Namespace)
			return nil
		}
	}

	err = deplClient.AddDeployment(db)
	if err != nil {
		return fmt.Errorf("Create DB failed: %+v", err)
	}

	return nil
}

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
			logrus.Infof("adding client: %+v", newClient)
		},
		func(oldClient types.Client, newClient types.Client) {
			logrus.Infof("updating client: %+v -> %+v", oldClient, newClient)
		},
		func(oldClient types.Client) {
			logrus.Infof("deleting client: %+v", oldClient)
		},
	)

	_, dbStopper := crdClient.DBListen(
		func(newDB types.Database) {
			logrus.Infof("adding database: %+v", newDB)

			err := addDeployment(deplClient, newDB)
			if err != nil {
				logrus.Errorf("Adding DB failed: %+v", err)
				return
			}

			logrus.Infof("database %s (%s) added", newDB.Name, newDB.Namespace)
		},
		func(oldDB types.Database, newDB types.Database) {
			logrus.Infof("updating database: %+v -> %+v", oldDB, newDB)
		},
		func(oldDB types.Database) {
			logrus.Infof("deleteting database %s (%s)", oldDB.Name, oldDB.Namespace)

			deleteDeployment(deplClient, oldDB)

			logrus.Infof("database %s (%s) deleted", oldDB.Name, oldDB.Namespace)
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

	log.Println("Stopped")
}
