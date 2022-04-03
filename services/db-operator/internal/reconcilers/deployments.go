package reconcilers

import (
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
)

type DeploymentReconciler struct {
	crdClient    *crds.DBClient
	deplClient   *deployments.DeploymentsClient
	databases    cache.Store
	statefulSets deployments.StatefulSetStore
	services     deployments.ServiceStore
}

func NewDeploymentReconciler(
	crdClient *crds.DBClient,
	deplClient *deployments.DeploymentsClient,
	databases cache.Store,
	statefulSets deployments.StatefulSetStore,
	services deployments.ServiceStore,
) *DeploymentReconciler {
	return &DeploymentReconciler{
		crdClient:    crdClient,
		deplClient:   deplClient,
		databases:    databases,
		statefulSets: statefulSets,
		services:     services,
	}
}

func (r *DeploymentReconciler) Start(
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
				logrus.Infof("stopped deployment reconciler")
			case <-events:
				timer = time.After(1 * time.Second)
			case <-timer:
				timer = time.After(60 * time.Second)

				requestedStatefulSets, requestedServices := requestedDatabases(r.databases)

				statefulSetsToAdd, statefulSetsToDelete := processStatefulSets(requestedStatefulSets, r.statefulSets)
				r.applyStatefulSets(statefulSetsToAdd, statefulSetsToDelete)

				servicesToAdd, servicesToDelete := processServices(requestedServices, r.services)
				r.applyServices(servicesToAdd, servicesToDelete)

				dbUpdates := r.databaseStatusUpdates()
				r.applyDatabaseUpdates(dbUpdates)
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
		if _, ok := requested[key]; !ok {
			item, _ := actual.GetByKey(key)
			toDelete[key] = item
		}
	}

	return toAdd, toDelete
}

func (r *DeploymentReconciler) applyStatefulSets(toAdd map[string]deployments.StatefulSet, toDelete map[string]deployments.StatefulSet) {
	for _, set := range toDelete {
		if err := r.deplClient.DeleteStatefulSet(set); err != nil {
			logrus.Errorf("Failed to delete stateful set: %+v", err)
		}
	}

	for _, set := range toAdd {
		if err := r.deplClient.AddStatefulSet(set); err != nil {
			logrus.Errorf("Failed to create stateful set: %+v", err)
		}
	}
}

func processServices(
	requested map[string]deployments.Service,
	actual deployments.ServiceStore,
) (map[string]deployments.Service, map[string]deployments.Service) {
	toAdd := map[string]deployments.Service{}
	toDelete := map[string]deployments.Service{}

	for key, requestedSet := range requested {
		setNew := false

		_, ok := actual.GetByKey(key)

		if !ok {
			setNew = true
		}

		if setNew {
			toAdd[key] = requestedSet
		}
	}

	for _, key := range actual.ListKeys() {
		if _, ok := requested[key]; !ok {
			item, _ := actual.GetByKey(key)
			toDelete[key] = item
		}
	}

	return toAdd, toDelete
}

func (r *DeploymentReconciler) applyServices(toAdd map[string]deployments.Service, toDelete map[string]deployments.Service) {
	for _, service := range toDelete {
		if err := r.deplClient.DeleteService(service); err != nil {
			logrus.Errorf("Failed to delete service: %+v", err)
		}
	}

	for _, service := range toAdd {
		if err := r.deplClient.AddService(service); err != nil {
			logrus.Errorf("Failed to create service: %+v", err)
		}
	}
}

func (r *DeploymentReconciler) databaseStatusUpdates() []deployments.StatefulSet {
	updates := []deployments.StatefulSet{}

	for _, key := range r.statefulSets.ListKeys() {
		set, _ := r.statefulSets.GetByKey(key)
		dbObj, exists, err := r.databases.GetByKey(key)

		if err != nil {
			logrus.Errorf("Error getting CockroachDB for %s, while checking for status updates", key)
			continue
		}

		if !exists {
			logrus.Warnf("stateful set without CockroachDB: %s", key)
			continue
		}

		db, ok := dbObj.(*crds.CockroachDB)
		if !ok {
			logrus.Errorf("Error converting %T into *CockroachDB", db)
			continue
		}

		if db.Status.Ready != set.Ready {
			updates = append(updates, set)
		}
	}

	return updates
}

func (r *DeploymentReconciler) applyDatabaseUpdates(updates []deployments.StatefulSet) {
	for _, db := range updates {
		if err := r.crdClient.DBUpdate(db.Name, db.Namespace, db.Ready); err != nil {
			logrus.Errorf("Failed updating cockroachdb CRD status: %+v", err)
		}
	}
}
