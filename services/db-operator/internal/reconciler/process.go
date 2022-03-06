package reconciler

// import (
// 	"fmt"

// 	"ponglehub.co.uk/operators/db/internal/actions"
// 	"ponglehub.co.uk/operators/db/internal/database"
// 	"ponglehub.co.uk/operators/db/internal/deployments"
// )

// func needsUpdate(request deployments.StatefulSet, actual deployments.StatefulSet) bool {
// 	return request.Name != actual.Name ||
// 		request.Namespace != actual.Namespace ||
// 		request.Storage != actual.Storage
// }

// func (r *Reconciler) processStatefulSets(
// 	actionChannel chan<- interface{},
// 	requested map[string]deployments.StatefulSet,
// 	actual map[string]deployments.StatefulSet,
// 	pending map[string]action,
// ) {
// 	for requestKey, requestedSet := range requested {
// 		setNew := false

// 		if pendingSet, ok := pending[fmt.Sprintf("%d:%T:%s", SET_ACTION, requestedSet, requestedSet.Key())]; ok {
// 			if needsUpdate(requestedSet, pendingSet.obj.(deployments.StatefulSet)) {
// 				setNew = true
// 				actionChannel <- actions.DeleteStatefulSet{StatefulSet: pendingSet.obj.(deployments.StatefulSet)}
// 			}
// 		} else if actualSet, ok := actual[requestKey]; ok {
// 			if needsUpdate(requestedSet, actualSet) {
// 				setNew = true
// 				actionChannel <- actions.DeleteStatefulSet{StatefulSet: actualSet}
// 			}
// 		} else {
// 			setNew = true
// 		}

// 		if setNew {
// 			actionChannel <- actions.AddStatefulSet{StatefulSet: requestedSet}
// 			r.removeAction(DELETE_ACTION, requestedSet)
// 			r.addAction(SET_ACTION, requestedSet)
// 		}
// 	}

// 	for actualKey, actualSet := range actual {
// 		_, deleting := pending[fmt.Sprintf("%d:%T:%s", DELETE_ACTION, actualSet, actualSet.Key())]
// 		_, requesting := requested[actualKey]

// 		if !requesting && !deleting {
// 			actionChannel <- actions.DeleteStatefulSet{StatefulSet: actualSet}
// 			r.removeAction(SET_ACTION, actualSet)
// 			r.addAction(DELETE_ACTION, actualSet)
// 		}
// 	}

// 	for _, action := range pending {
// 		if action.code != SET_ACTION {
// 			continue
// 		}

// 		set, isSet := action.obj.(deployments.StatefulSet)
// 		if !isSet {
// 			continue
// 		}

// 		_, requesting := requested[set.Key()]
// 		if requesting {
// 			continue
// 		}

// 		actionChannel <- actions.DeleteStatefulSet{StatefulSet: set}
// 		r.removeAction(SET_ACTION, set)
// 		r.addAction(DELETE_ACTION, set)
// 	}
// }

// func (r *Reconciler) processServices(
// 	actionChannel chan<- interface{},
// 	requested map[string]deployments.Service,
// 	actual map[string]deployments.Service,
// 	pending map[string]action,
// ) {

// 	for requestKey, requestedService := range requested {
// 		setNew := false

// 		if pendingService, ok := pending[fmt.Sprintf("%d:%T:%s", SET_ACTION, requestedService, requestedService.Key())]; ok {
// 			if requestedService != pendingService.obj {
// 				setNew = true
// 				actionChannel <- actions.DeleteService{Service: pendingService.obj.(deployments.Service)}
// 			}
// 		} else if actualService, ok := actual[requestKey]; ok {
// 			if requestedService != actualService {
// 				setNew = true
// 				actionChannel <- actions.DeleteService{Service: actualService}
// 			}
// 		} else {
// 			setNew = true
// 		}

// 		if setNew {
// 			actionChannel <- actions.AddService{Service: requestedService}
// 			r.removeAction(DELETE_ACTION, requestedService)
// 			r.addAction(SET_ACTION, requestedService)
// 		}
// 	}

// 	for actualKey, actualService := range actual {
// 		_, deleting := pending[fmt.Sprintf("%d:%T:%s", DELETE_ACTION, actualService, actualService.Key())]
// 		_, requesting := requested[actualKey]

// 		if !requesting && !deleting {
// 			actionChannel <- actions.DeleteService{Service: actualService}
// 			r.removeAction(SET_ACTION, actualService)
// 			r.addAction(DELETE_ACTION, actualService)
// 		}
// 	}

// 	for _, action := range pending {
// 		if action.code != SET_ACTION {
// 			continue
// 		}

// 		service, isService := action.obj.(deployments.Service)
// 		if !isService {
// 			continue
// 		}

// 		_, requesting := requested[service.Key()]
// 		if requesting {
// 			continue
// 		}

// 		actionChannel <- actions.DeleteService{Service: service}
// 		r.removeAction(SET_ACTION, service)
// 		r.addAction(DELETE_ACTION, service)
// 	}
// }

// func (r *Reconciler) processClients(
// 	actionChannel chan<- interface{},
// 	requested map[string]database.Client,
// 	actual map[string]database.Client,
// 	pending map[string]action,
// ) {
// }

// func (r *Reconciler) processSecrets(
// 	actionChannel chan<- interface{},
// 	requested map[string]deployments.ClientSecret,
// 	actual map[string]deployments.ClientSecret,
// 	pending map[string]action,
// ) {
// }
