package database

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/lib/postgres/pkg/connect"
	"ponglehub.co.uk/lib/postgres/pkg/migrate"
)

type DatabaseClient struct {
	deployments map[string]Client
}

func New() *DatabaseClient {
	return &DatabaseClient{
		deployments: map[string]Client{},
	}
}

func (d *DatabaseClient) CreateClient(client Client) error {
	logrus.Infof("Creating a new database client: %v", client)

	config := connect.ConnectConfig{
		Host:     fmt.Sprintf("%s.%s.svc.cluster.local", client.Deployment, client.Namespace),
		Port:     26257,
		Username: "root",
		Database: client.Database,
	}

	err := migrate.Initialize(config, client.Database, client.Username)
	if err != nil {
		return fmt.Errorf("failed to initialise user or database: %+v", err)
	}

	d.deployments[client.Key()] = client
	return nil
}

func (d *DatabaseClient) DeleteClient(client Client) error {
	logrus.Infof("Removing a database client: %v", client)

	config := connect.ConnectConfig{
		Host:     fmt.Sprintf("%s.%s.svc.cluster.local", client.Deployment, client.Namespace),
		Port:     26257,
		Username: "root",
		Database: client.Database,
	}

	err := migrate.UnInitialize(config, client.Database, client.Username)
	if err != nil {
		return fmt.Errorf("failed to initialise user or database: %+v", err)
	}

	delete(d.deployments, client.Key())
	return nil
}

func (d *DatabaseClient) HasClient(client Client) bool {
	_, ok := d.deployments[client.Key()]
	return ok
}

func (d *DatabaseClient) ListClients() map[string]Client {
	return d.deployments
}

func (d *DatabaseClient) PruneClient(client Client) {
	delete(d.deployments, client.Key())
}
