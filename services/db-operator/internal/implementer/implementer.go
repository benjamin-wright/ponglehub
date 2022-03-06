package implementer

import (
	"github.com/sirupsen/logrus"
	"ponglehub.co.uk/operators/db/internal/crds"
	"ponglehub.co.uk/operators/db/internal/deployments"
)

type Implementer struct {
	crdClient  *crds.Client
	deplClient *deployments.DeploymentsClient
}

func New(crdClient *crds.Client, deplClient *deployments.DeploymentsClient) *Implementer {
	return &Implementer{
		crdClient:  crdClient,
		deplClient: deplClient,
	}
}

func (i *Implementer) DeleteStatefulSets(sets map[string]deployments.StatefulSet) {
	for _, set := range sets {
		err := i.deplClient.DeleteStatefulSet(set)
		if err != nil {
			logrus.Errorf("Failed to delete stateful set: %s (%s)", set.Name, set.Namespace)
		}
	}
}

func (i *Implementer) AddStatefulSets(sets map[string]deployments.StatefulSet) {
	for _, set := range sets {
		err := i.deplClient.AddStatefulSet(set)
		if err != nil {
			logrus.Errorf("Failed to delete stateful set: %s (%s)", set.Name, set.Namespace)
		}
	}
}
