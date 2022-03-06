package deployments

import (
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"
)

type StatefulSetStore struct {
	store cache.Store
}

func (s StatefulSetStore) ListKeys() []string {
	return s.store.ListKeys()
}

func (s StatefulSetStore) GetByKey(key string) (StatefulSet, bool) {
	item, exists, err := s.store.GetByKey(key)

	if err != nil {
		logrus.Errorf("Failed to fetch \"%s\" from stateful set store", key)
		return StatefulSet{}, false
	}

	if !exists {
		return StatefulSet{}, false
	}

	kubeObj, ok := item.(*appsv1.StatefulSet)
	if !ok {
		logrus.Errorf("Failed to convert stateful set pointer from store: %T", item)
		return StatefulSet{}, false
	}

	set, err := fromSS(kubeObj)
	if err != nil {
		logrus.Errorf("Failed to parse \"%s\" from stateful set store: %+v", key, err)
		return StatefulSet{}, false
	}

	return set, true
}
