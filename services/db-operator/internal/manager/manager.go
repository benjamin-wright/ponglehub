package manager

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/operators/db/internal/certs"
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

	if err := deplClient.DeleteNodeSecret(db.Namespace, db.Name); err != nil {
		logrus.Errorf("failed to delete node secret: %+v", err)
	}

	if err := deplClient.DeleteCASecret(db.Namespace, db.Name); err != nil {
		logrus.Errorf("failed to delete ca secret: %+v", err)
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
		return fmt.Errorf("create DB failed: %+v", err)
	}

	return nil
}
