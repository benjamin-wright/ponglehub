package reconciler

import (
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
)

type Reconciler struct {
	crdClient    *crds.DBClient
	deplClient   *deployments.DeploymentsClient
	databases    cache.Store
	clients      cache.Store
	statefulSets deployments.StatefulSetStore
	services     deployments.ServiceStore
	secrets      cache.Store
}

func New(
	crdClient *crds.DBClient,
	deplClient *deployments.DeploymentsClient,
	databases cache.Store,
	clients cache.Store,
	statefulSets deployments.StatefulSetStore,
	services deployments.ServiceStore,
	secrets cache.Store,
) *Reconciler {
	return &Reconciler{
		crdClient:    crdClient,
		deplClient:   deplClient,
		databases:    databases,
		clients:      clients,
		statefulSets: statefulSets,
		services:     services,
		secrets:      secrets,
	}
}

func (r *Reconciler) Start(
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
				logrus.Infof("stopped reconciler")
			case <-events:
				logrus.Infof("update")
				timer = time.After(5 * time.Second)
			case <-timer:
				logrus.Infof("reconciling")
				requestedStatefulSets, _ := requestedDatabases(r.databases)

				statefulSetsToAdd, statefulSetsToDelete := processStatefulSets(requestedStatefulSets, r.statefulSets)

				for _, set := range statefulSetsToDelete {
					if err := r.deplClient.DeleteStatefulSet(set); err != nil {
						logrus.Errorf("Failed to delete stateful set: %+v", err)
					}
				}

				for _, set := range statefulSetsToAdd {
					if err := r.deplClient.AddStatefulSet(set); err != nil {
						logrus.Errorf("Failed to create stateful set: %+v", err)
					}
				}
			}
		}
	}(stopper)

	return stopper
}

func requestedDatabases(databases cache.Store) (map[string]deployments.StatefulSet, map[string]deployments.Service) {
	sets := map[string]deployments.StatefulSet{}
	services := map[string]deployments.Service{}

	for _, database := range databases.List() {
		db, ok := database.(*crds.CockroachDB)
		if !ok {
			logrus.Warnf("Failed to convert database: %T", db)
			continue
		}

		key := db.Namespace + "/" + db.Name

		sets[key] = deployments.StatefulSet{
			Name:      db.Name,
			Namespace: db.Namespace,
			Storage:   db.Spec.Storage,
		}

		services[key] = deployments.Service{
			Name:      db.Name,
			Namespace: db.Namespace,
		}
	}

	return sets, services
}

func processStatefulSets(
	requested map[string]deployments.StatefulSet,
	actual deployments.StatefulSetStore,
) (map[string]deployments.StatefulSet, map[string]deployments.StatefulSet) {
	toAdd := map[string]deployments.StatefulSet{}
	toDelete := map[string]deployments.StatefulSet{}

	for key, requestedSet := range requested {
		setNew := false

		actualSet, ok := actual.GetByKey(key)

		if ok {
			if actualSet.Storage != requestedSet.Storage {
				setNew = true
				toDelete[key] = actualSet
			}
		} else {
			setNew = true
		}

		if setNew {
			toAdd[key] = requestedSet
		}
	}

	for _, key := range actual.ListKeys() {
		if item, ok := requested[key]; !ok {
			toDelete[key] = item
		}
	}

	return toAdd, toDelete
}

// func New() *Reconciler {
// 	return &Reconciler{
// 		requestedDatabases: map[string]crds.Database{},
// 		requestedClients:   map[string]crds.Client{},
// 		actualStatefulSets: map[string]deployments.StatefulSet{},
// 		actualServices:     map[string]deployments.Service{},
// 		actualSecrets:      map[string]deployments.ClientSecret{},
// 		actualClients:      map[string]database.Client{},
// 		pendingActions:     map[string]action{},
// 	}
// }

// func (r *Reconciler) SetDatabase(database crds.Database) bool {
// 	key := database.Key()

// 	if existing, ok := r.requestedDatabases[key]; ok {
// 		if existing == database {
// 			return false
// 		}
// 	}

// 	r.requestedDatabases[key] = database
// 	return true
// }

// func (r *Reconciler) RemoveDatabase(database crds.Database) bool {
// 	key := database.Key()

// 	if _, ok := r.requestedDatabases[key]; ok {
// 		delete(r.requestedDatabases, key)
// 		return true
// 	}

// 	return false
// }

// func (r *Reconciler) SetClient(client crds.Client) bool {
// 	key := client.Key()

// 	if existing, ok := r.requestedClients[key]; ok {
// 		if existing == client {
// 			return false
// 		}
// 	}

// 	r.requestedClients[key] = client
// 	return true
// }

// func (r *Reconciler) RemoveClient(client crds.Client) bool {
// 	key := client.Key()

// 	if _, ok := r.requestedClients[key]; ok {
// 		delete(r.requestedClients, key)
// 		return true
// 	}

// 	return false
// }

// func (r *Reconciler) SetStatefulSet(statefulset deployments.StatefulSet) bool {
// 	key := statefulset.Key()

// 	if existing, ok := r.actualStatefulSets[key]; ok {
// 		if existing == statefulset {
// 			return false
// 		}
// 	}

// 	r.actualStatefulSets[key] = statefulset
// 	r.removeAction(SET_ACTION, statefulset)
// 	return true
// }

// func (r *Reconciler) RemoveStatefulSet(statefulset deployments.StatefulSet) bool {
// 	key := statefulset.Key()

// 	if _, ok := r.actualStatefulSets[key]; ok {
// 		delete(r.actualStatefulSets, key)
// 		r.removeAction(DELETE_ACTION, statefulset)
// 		return true
// 	}

// 	return false
// }

// func (r *Reconciler) SetService(service deployments.Service) bool {
// 	key := service.Key()

// 	if existing, ok := r.actualServices[key]; ok {
// 		if existing == service {
// 			return false
// 		}
// 	}

// 	r.actualServices[key] = service
// 	return true
// }

// func (r *Reconciler) RemoveService(service deployments.Service) bool {
// 	key := service.Key()

// 	if _, ok := r.actualServices[key]; ok {
// 		delete(r.actualServices, key)
// 		return true
// 	}

// 	return false
// }

// func (r *Reconciler) SetClientSecret(client deployments.ClientSecret) bool {
// 	key := client.Key()

// 	if existing, ok := r.actualSecrets[key]; ok {
// 		if existing == client {
// 			return false
// 		}
// 	}

// 	r.actualSecrets[key] = client
// 	return true
// }

// func (r *Reconciler) RemoveClientSecret(client deployments.ClientSecret) bool {
// 	key := client.Key()

// 	if _, ok := r.actualSecrets[key]; ok {
// 		delete(r.actualSecrets, key)
// 		return true
// 	}

// 	return false
// }

// func (r *Reconciler) addAction(code actionCode, object keyable) {
// 	key := fmt.Sprintf("%d:%T:%s", code, object, object.Key())
// 	r.pendingActions[key] = action{code, object}
// }

// func (r *Reconciler) removeAction(code actionCode, object keyable) {
// 	key := fmt.Sprintf("%d:%T:%s", code, object, object.Key())
// 	if actual, ok := r.pendingActions[key]; ok && actual.obj == object {
// 		delete(r.pendingActions, key)
// 	}
// }

// func (r *Reconciler) Reconcile(actions chan<- interface{}) {
// 	requestedStatefulSets, requestedServices := getDatabaseRequests(r.requestedDatabases)
// 	requestedClients := getClientRequests(r.requestedClients, r.actualStatefulSets)
// 	requestedSecrets := getSecretRequests(r.requestedClients, r.actualClients)

// 	r.processStatefulSets(actions, requestedStatefulSets, r.actualStatefulSets, r.pendingActions)
// 	r.processServices(actions, requestedServices, r.actualServices, r.pendingActions)
// 	r.processClients(actions, requestedClients, r.actualClients, r.pendingActions)
// 	r.processSecrets(actions, requestedSecrets, r.actualSecrets, r.pendingActions)
// }
