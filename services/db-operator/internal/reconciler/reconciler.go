package reconciler

import (
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
)

type Reconciler struct {
	requestedDatabases map[string]crds.Database
	requestedClients   map[string]crds.Client
	actualStatefulSets map[string]deployments.StatefulSet
	actualServices     map[string]deployments.Service
}

func (r *Reconciler) SetDatabase(database crds.Database) bool {
	key := database.Key()

	if existing, ok := r.requestedDatabases[key]; ok {
		if existing == database {
			return false
		}
	}

	r.requestedDatabases[key] = database
	return true
}

func (r *Reconciler) RemoveDatabase(database crds.Database) bool {
	key := database.Key()

	if _, ok := r.requestedDatabases[key]; ok {
		delete(r.requestedDatabases, key)
		return true
	}

	return false
}

func (r *Reconciler) SetClient(client crds.Client) bool {
	key := client.Key()

	if existing, ok := r.requestedClients[key]; ok {
		if existing == client {
			return false
		}
	}

	r.requestedClients[key] = client
	return true
}

func (r *Reconciler) RemoveClient(client crds.Client) bool {
	key := client.Key()

	if _, ok := r.requestedClients[key]; ok {
		delete(r.requestedClients, key)
		return true
	}

	return false
}

func (r *Reconciler) SetStatefulSet(statefulset deployments.StatefulSet) bool {
	key := statefulset.Key()

	if existing, ok := r.actualStatefulSets[key]; ok {
		if existing == statefulset {
			return false
		}
	}

	r.actualStatefulSets[key] = statefulset
	return true
}

func (r *Reconciler) RemoveStatefulSet(client deployments.StatefulSet) bool {
	key := client.Key()

	if _, ok := r.actualStatefulSets[key]; ok {
		delete(r.actualStatefulSets, key)
		return true
	}

	return false
}

func (r *Reconciler) SetService(service deployments.Service) bool {
	key := service.Key()

	if existing, ok := r.actualServices[key]; ok {
		if existing == service {
			return false
		}
	}

	r.actualServices[key] = service
	return true
}

func (r *Reconciler) RemoveService(client deployments.Service) bool {
	key := client.Key()

	if _, ok := r.actualServices[key]; ok {
		delete(r.actualServices, key)
		return true
	}

	return false
}

func (r *Reconciler) Reconcile(actions chan<- interface{}) {
	requestedSets, requestedServices := getRequested(r.requestedDatabases)
}

func getRequested(databases map[string]crds.Database) (map[string]deployments.StatefulSet, map[string]deployments.Service) {
	sets := map[string]deployments.StatefulSet{}
	services := map[string]deployments.Service{}

	for _, db := range databases {
		set := deployments.StatefulSet{
			Name:      db.Name,
			Namespace: db.Namespace,
			Storage:   db.Storage,
		}
		setKey := set.Key()
		sets[setKey] = set

		service := deployments.Service{
			Name:      db.Name,
			Namespace: db.Namespace,
		}
		serviceKey := service.Key()
		services[serviceKey] = service
	}

	return sets, services
}
