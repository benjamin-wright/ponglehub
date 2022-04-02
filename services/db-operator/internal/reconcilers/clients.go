package reconcilers

import (
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/database"
	"ponglehub.co.uk/operators/db/internal/deployments"
)

type ClientReconciler struct {
	crdClient    *crds.DBClient
	deplClient   *deployments.DeploymentsClient
	databases    cache.Store
	clients      cache.Store
	statefulSets deployments.StatefulSetStore
}

func NewClientReconciler(
	crdClient *crds.DBClient,
	deplClient *deployments.DeploymentsClient,
	databases cache.Store,
	clients cache.Store,
	statefulSets deployments.StatefulSetStore,
) *ClientReconciler {
	return &ClientReconciler{
		crdClient:    crdClient,
		deplClient:   deplClient,
		databases:    databases,
		clients:      clients,
		statefulSets: statefulSets,
	}
}

func (r *ClientReconciler) Start(
	events <-chan interface{},
) chan<- interface{} {

	stopper := make(chan interface{})

	go func(stopper <-chan interface{}) {
		running := true
		timer := time.After(5 * time.Second)

		for running {
			select {
			case <-stopper:
				running = false
				logrus.Infof("stopped client reconciler")
			case <-events:
				timer = time.After(1 * time.Second)
			case <-timer:
				timer = time.After(60 * time.Second)

				requestedClients := requestedClients(r.databases, r.clients)
				clientsToAdd, clientsToRemove := processClients(requestedClients, r.clients)
			}
		}
	}(stopper)

	return stopper
}

func requestedClients(databases cache.Store, clients cache.Store) map[string]database.Client {
	requests := map[string]database.Client{}

	for _, client := range clients.List() {
		cli, ok := client.(*crds.CockroachClient)
		if !ok {
			logrus.Warnf("Failed to convert client: %T", cli)
			continue
		}

		found := false
		for _, database := range databases.List() {
			db, ok := database.(*crds.CockroachDB)
			if !ok {
				logrus.Warnf("Failed to convert database: %T", db)
				continue
			}

			if cli.Namespace == db.Namespace && cli.Spec.Deployment == db.Name {
				found = true
				break
			}
		}

		if !found {
			continue
		}

		key := cli.Namespace + "/" + cli.Name
		requests[key] = database.Client{
			Username:   cli.Spec.Username,
			Deployment: cli.Spec.Deployment,
			Database:   cli.Spec.Database,
		}
	}

	return requests
}

func processClients(requested map[string]database.Client, clients cache.Store) (map[string]database.Client, map[string]database.Client) {
	return nil, nil
}
