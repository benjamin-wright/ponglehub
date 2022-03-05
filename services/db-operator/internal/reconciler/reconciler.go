package reconciler

import (
	"fmt"

	"ponglehub.co.uk/operators/db/internal/actions"
	"ponglehub.co.uk/operators/db/internal/crds"
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
	pendingActions     map[string]action
}

func New() *Reconciler {
	return &Reconciler{
		requestedDatabases: map[string]crds.Database{},
		requestedClients:   map[string]crds.Client{},
		actualStatefulSets: map[string]deployments.StatefulSet{},
		actualServices:     map[string]deployments.Service{},
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

func (r *Reconciler) RemoveService(client deployments.Service) bool {
	key := client.Key()

	if _, ok := r.actualServices[key]; ok {
		delete(r.actualServices, key)
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
	delete(r.pendingActions, key)
}

func (r *Reconciler) Reconcile(actions chan<- interface{}) {
	requestedStatefulSets, requestedServices := getRequested(r.requestedDatabases)

	r.processStatefulSets(actions, requestedStatefulSets, r.actualStatefulSets, r.pendingActions)
	r.processServices(actions, requestedServices, r.actualServices, r.pendingActions)
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

func (r *Reconciler) processStatefulSets(
	actionChannel chan<- interface{},
	requested map[string]deployments.StatefulSet,
	actual map[string]deployments.StatefulSet,
	pending map[string]action,
) {
	for requestKey, requestedSet := range requested {
		setNew := false

		if pendingSet, ok := pending[fmt.Sprintf("%d:%T:%s", SET_ACTION, requestedSet, requestedSet.Key())]; ok {
			if requestedSet != pendingSet.obj {
				setNew = true
				actionChannel <- actions.DeleteStatefulSet{StatefulSet: pendingSet.obj.(deployments.StatefulSet)}
			}
		} else if actualSet, ok := actual[requestKey]; ok {
			if requestedSet != actualSet {
				setNew = true
				actionChannel <- actions.DeleteStatefulSet{StatefulSet: actualSet}
			}
		} else {
			setNew = true
		}

		if setNew {
			actionChannel <- actions.AddStatefulSet{StatefulSet: requestedSet}
			r.removeAction(DELETE_ACTION, requestedSet)
			r.addAction(SET_ACTION, requestedSet)
		}
	}

	for actualKey, actualSet := range actual {
		_, deleting := pending[fmt.Sprintf("%d:%T:%s", DELETE_ACTION, actualSet, actualSet.Key())]
		_, requesting := requested[actualKey]

		if !requesting && !deleting {
			actionChannel <- actions.DeleteStatefulSet{StatefulSet: actualSet}
			r.removeAction(SET_ACTION, actualSet)
			r.addAction(DELETE_ACTION, actualSet)
		}
	}

	for _, action := range pending {
		if action.code != SET_ACTION {
			continue
		}

		set, isSet := action.obj.(deployments.StatefulSet)
		if !isSet {
			continue
		}

		_, requesting := requested[set.Key()]
		if requesting {
			continue
		}

		actionChannel <- actions.DeleteStatefulSet{StatefulSet: set}
		r.removeAction(SET_ACTION, set)
		r.addAction(DELETE_ACTION, set)
	}
}

func (r *Reconciler) processServices(
	actionChannel chan<- interface{},
	requested map[string]deployments.Service,
	actual map[string]deployments.Service,
	pending map[string]action,
) {

	for requestKey, requestedService := range requested {
		setNew := false

		if pendingService, ok := pending[fmt.Sprintf("%d:%T:%s", SET_ACTION, requestedService, requestedService.Key())]; ok {
			if requestedService != pendingService.obj {
				setNew = true
				actionChannel <- actions.DeleteService{Service: pendingService.obj.(deployments.Service)}
			}
		} else if actualService, ok := actual[requestKey]; ok {
			if requestedService != actualService {
				setNew = true
				actionChannel <- actions.DeleteService{Service: actualService}
			}
		} else {
			setNew = true
		}

		if setNew {
			actionChannel <- actions.AddService{Service: requestedService}
			r.removeAction(DELETE_ACTION, requestedService)
			r.addAction(SET_ACTION, requestedService)
		}
	}

	for actualKey, actualService := range actual {
		_, deleting := pending[fmt.Sprintf("%d:%T:%s", DELETE_ACTION, actualService, actualService.Key())]
		_, requesting := requested[actualKey]

		if !requesting && !deleting {
			actionChannel <- actions.DeleteService{Service: actualService}
			r.removeAction(SET_ACTION, actualService)
			r.addAction(DELETE_ACTION, actualService)
		}
	}

	for _, action := range pending {
		if action.code != SET_ACTION {
			continue
		}

		service, isService := action.obj.(deployments.Service)
		if !isService {
			continue
		}

		_, requesting := requested[service.Key()]
		if requesting {
			continue
		}

		actionChannel <- actions.DeleteService{Service: service}
		r.removeAction(SET_ACTION, service)
		r.addAction(DELETE_ACTION, service)
	}
}
