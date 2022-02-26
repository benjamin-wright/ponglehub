package reconciler

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
)

func Start(events <-chan interface{}, actions chan<- interface{}) chan<- interface{} {
	stopper := make(chan interface{})
	reconciler := &Reconciler{}

	go func() {
		running := true
		for running {
			changed := false

			select {
			case <-stopper:
				running = false
			case e := <-events:
				logrus.Infof("EVENT: %T", e)

				switch event := e.(type) {
				case crds.ClientAddedEvent:
					changed = reconciler.SetClient(event.New)
				case crds.ClientUpdatedEvent:
					changed = reconciler.SetClient(event.New)
				case crds.ClientDeletedEvent:
					changed = reconciler.RemoveClient(event.Old)
				case crds.DatabaseAddedEvent:
					changed = reconciler.SetDatabase(event.New)
				case crds.DatabaseUpdatedEvent:
					changed = reconciler.SetDatabase(event.New)
				case crds.DatabaseDeletedEvent:
					changed = reconciler.RemoveDatabase(event.Old)
				case deployments.ServiceAddedEvent:
					changed = reconciler.SetService(event.New)
				case deployments.ServiceUpdatedEvent:
					changed = reconciler.SetService(event.New)
				case deployments.ServiceDeletedEvent:
					changed = reconciler.RemoveService(event.Old)
				case deployments.StatefulSetAddedEvent:
					changed = reconciler.SetStatefulSet(event.New)
				case deployments.StatefulSetUpdatedEvent:
					changed = reconciler.SetStatefulSet(event.New)
				case deployments.StatefulSetDeletedEvent:
					changed = reconciler.RemoveStatefulSet(event.Old)
				default:
					logrus.Errorf("unrecognised event: %+v", event)
				}
			}

			if changed {
				logrus.Infof("Something has changed!")
				reconciler.Reconcile(actions)
			}
		}
	}()

	return stopper
}
