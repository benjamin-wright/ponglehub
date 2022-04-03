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
	dbClient     *database.DatabaseClient
	databases    cache.Store
	clients      cache.Store
	statefulSets deployments.StatefulSetStore
}

func NewClientReconciler(
	crdClient *crds.DBClient,
	deplClient *deployments.DeploymentsClient,
	dbClient *database.DatabaseClient,
	databases cache.Store,
	clients cache.Store,
	statefulSets deployments.StatefulSetStore,
) *ClientReconciler {
	return &ClientReconciler{
		crdClient:    crdClient,
		deplClient:   deplClient,
		dbClient:     dbClient,
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

				requestedClients := requestedClients(r.databases, r.clients, r.dbClient)
				clientsToAdd, clientsToRemove := processClients(requestedClients, r.dbClient)
				r.applyClients(clientsToAdd, clientsToRemove)

				clientUpdates := r.clientStatusUpdates()
				r.applyClientUpdates(clientUpdates)
			}
		}
	}(stopper)

	return stopper
}

func requestedClients(databases cache.Store, clients cache.Store, dbClient *database.DatabaseClient) map[string]database.Client {
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

			if cli.Namespace == db.Namespace && cli.Spec.Deployment == db.Name && db.Status.Ready {
				found = true
				break
			}
		}

		newClient := database.Client{
			Username:   cli.Spec.Username,
			Deployment: cli.Spec.Deployment,
			Database:   cli.Spec.Database,
			Namespace:  cli.Namespace,
		}

		if !found {
			if dbClient.HasClient(newClient) {
				dbClient.PruneClient(newClient)
			}

			continue
		}

		requests[newClient.Key()] = newClient
	}

	return requests
}

func processClients(requested map[string]database.Client, clients *database.DatabaseClient) (map[string]database.Client, map[string]database.Client) {
	toAdd := map[string]database.Client{}
	toRemove := map[string]database.Client{}

	for key, client := range requested {
		if !clients.HasClient(client) {
			toAdd[key] = client
		}
	}

	for key, existing := range clients.ListClients() {
		missing := true
		for _, client := range requested {
			if client == existing {
				missing = false
				break
			}
		}

		if missing {
			toRemove[key] = existing
		}
	}

	return toAdd, toRemove
}

func (r *ClientReconciler) applyClients(toAdd map[string]database.Client, toRemove map[string]database.Client) {
	for _, client := range toAdd {
		r.dbClient.CreateClient(client)
	}

	for _, client := range toRemove {
		r.dbClient.DeleteClient(client)
	}
}

type ClientUpdate struct {
	Name      string
	Namespace string
	Ready     bool
}

func (r *ClientReconciler) clientStatusUpdates() []ClientUpdate {
	updates := []ClientUpdate{}

	for _, client := range r.clients.List() {
		cli, ok := client.(*crds.CockroachClient)
		if !ok {
			logrus.Warnf("Failed to convert client: %T", cli)
			continue
		}

		requested := database.Client{
			Username:   cli.Spec.Username,
			Deployment: cli.Spec.Deployment,
			Database:   cli.Spec.Database,
			Namespace:  cli.Namespace,
		}

		exists := r.dbClient.HasClient(requested)

		if exists && !cli.Status.Ready {
			updates = append(updates, ClientUpdate{Name: cli.Name, Namespace: cli.Namespace, Ready: true})
		}

		if !exists && cli.Status.Ready {
			updates = append(updates, ClientUpdate{Name: cli.Name, Namespace: cli.Namespace, Ready: false})
		}
	}

	return updates
}

func (r *ClientReconciler) applyClientUpdates(updates []ClientUpdate) {
	for _, update := range updates {
		if err := r.crdClient.ClientUpdate(update.Name, update.Namespace, update.Ready); err != nil {
			logrus.Errorf("Failed updating cockroachclient CRD status: %+v", err)
		}
	}
}
