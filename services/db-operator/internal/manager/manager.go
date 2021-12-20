package manager

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
	"ponglehub.co.uk/lib/postgres/pkg/migrate"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
	"ponglehub.co.uk/operators/db/internal/types"
)

func DeleteDeployment(deplClient *deployments.DeploymentsClient, db types.Database) {
	if err := deplClient.DeleteDeployment(db); err != nil {
		logrus.Errorf("failed to delete database: %+v", err)
	}

	if err := deplClient.DeleteService(db.Namespace, db.Name); err != nil {
		logrus.Errorf("failed to delete service: %+v", err)
	}
}

func AddDeployment(deplClient *deployments.DeploymentsClient, db types.Database) error {
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
		return fmt.Errorf("create DB failed: %+v", err)
	}

	return nil
}

func DeleteClient(deplClient *deployments.DeploymentsClient, client types.Client) error {
	return nil
}

func generatePassword() string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials

	length := 16
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]

	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}

	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})

	return string(buf) // E.g. "3i[g0|)z"
}

func AddClient(crdClient *crds.DBClient, client types.Client) error {
	err := migrate.Initialize(
		connect.ConnectConfig{
			Host:     fmt.Sprintf("%s.%s.svc.cluster.local", client.Deployment, client.Namespace),
			Port:     26257,
			Username: "root",
			Database: client.Database,
		},
		client.Database,
		client.Username,
	)
	if err != nil {
		return fmt.Errorf("failed to initialise database %s for client %s (%s): %+v", client.Database, client.Name, client.Namespace, err)
	}

	err = crdClient.ClientUpdate(client.Name, client.Namespace, true)
	if err != nil {
		return fmt.Errorf("failed to update client secret for %s (%s): %+v", client.Name, client.Namespace, err)
	}
	err = crdClient.ClientUpdate(client.Name, client.Namespace, true)
	if err != nil {
		return fmt.Errorf("failed to update client secret for %s (%s): %+v", client.Name, client.Namespace, err)
	}
	err = crdClient.ClientUpdate(client.Name, client.Namespace, true)
	if err != nil {
		return fmt.Errorf("failed to update client secret for %s (%s): %+v", client.Name, client.Namespace, err)
	}

	return nil
}
