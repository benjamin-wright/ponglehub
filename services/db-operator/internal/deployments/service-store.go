package deployments

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

type ServiceStore struct {
	store cache.Store
}

func (s ServiceStore) ListKeys() []string {
	return s.store.ListKeys()
}

func (s ServiceStore) GetByKey(key string) (Service, bool) {
	item, exists, err := s.store.GetByKey(key)

	if err != nil {
		logrus.Errorf("Failed to fetch \"%s\" from service store", key)
		return Service{}, false
	}

	if !exists {
		return Service{}, false
	}

	kubeObj, ok := item.(*corev1.Service)
	if !ok {
		logrus.Errorf("Failed to convert service pointer from store: %T", item)
		return Service{}, false
	}

	service, err := fromService(kubeObj)
	if err != nil {
		logrus.Errorf("Failed to parse \"%s\" from service store: %+v", key, err)
		return Service{}, false
	}

	return service, true
}
