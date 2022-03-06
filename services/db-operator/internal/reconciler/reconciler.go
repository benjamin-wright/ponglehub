package reconciler

import (
	"fmt"

	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/database"
	"ponglehub.co.uk/operators/db/internal/deployments"
)

type actionCode int64

type keyable interface {
	Key() string
}

type action struct {
	code actionCode
	obj  keyable
}

const (
	SET_ACTION actionCode = iota
	DELETE_ACTION
)

type Reconciler struct {
	requestedDatabases map[string]crds.Database
	requestedClients   map[string]crds.Client
	actualStatefulSets map[string]deployments.StatefulSet
	actualServices     map[string]deployments.Service
	actualSecrets      map[string]deployments.ClientSecret
	actualClients      map[string]database.Client
	pendingActions     map[string]action
}

func New() *Reconciler {
	return &Reconciler{
		requestedDatabases: map[string]crds.Database{},
		requestedClients:   map[string]crds.Client{},
		actualStatefulSets: map[string]deployments.StatefulSet{},
		actualServices:     map[string]deployments.Service{},
		actualSecrets:      map[string]deployments.ClientSecret{},
		actualClients:      map[string]database.Client{},
		pendingActions:     map[string]action{},
	}
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
	r.removeAction(SET_ACTION, statefulset)
	return true
}

func (r *Reconciler) RemoveStatefulSet(statefulset deployments.StatefulSet) bool {
	key := statefulset.Key()

	if _, ok := r.actualStatefulSets[key]; ok {
		delete(r.actualStatefulSets, key)
		r.removeAction(DELETE_ACTION, statefulset)
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

func (r *Reconciler) RemoveService(service deployments.Service) bool {
	key := service.Key()

	if _, ok := r.actualServices[key]; ok {
		delete(r.actualServices, key)
		return true
	}

	return false
}

func (r *Reconciler) SetClientSecret(client deployments.ClientSecret) bool {
	key := client.Key()

	if existing, ok := r.actualSecrets[key]; ok {
		if existing == client {
			return false
		}
	}

	r.actualSecrets[key] = client
	return true
}

func (r *Reconciler) RemoveClientSecret(client deployments.ClientSecret) bool {
	key := client.Key()

	if _, ok := r.actualSecrets[key]; ok {
		delete(r.actualSecrets, key)
		return true
	}

	return false
}

func (r *Reconciler) addAction(code actionCode, object keyable) {
	key := fmt.Sprintf("%d:%T:%s", code, object, object.Key())
	r.pendingActions[key] = action{code, object}
}

func (r *Reconciler) removeAction(code actionCode, object keyable) {
	key := fmt.Sprintf("%d:%T:%s", code, object, object.Key())
	if actual, ok := r.pendingActions[key]; ok && actual.obj == object {
		delete(r.pendingActions, key)
	}
}

func (r *Reconciler) Reconcile(actions chan<- interface{}) {
	requestedStatefulSets, requestedServices := getDatabaseRequests(r.requestedDatabases)
	requestedClients := getClientRequests(r.requestedClients, r.actualStatefulSets)
	requestedSecrets := getSecretRequests(r.requestedClients, r.actualClients)

	r.processStatefulSets(actions, requestedStatefulSets, r.actualStatefulSets, r.pendingActions)
	r.processServices(actions, requestedServices, r.actualServices, r.pendingActions)
	r.processClients(actions, requestedClients, r.actualClients, r.pendingActions)
	r.processSecrets(actions, requestedSecrets, r.actualSecrets, r.pendingActions)
}
