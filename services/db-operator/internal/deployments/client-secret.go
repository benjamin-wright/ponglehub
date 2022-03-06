package deployments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ClientSecret struct {
	Namespace string
	Name      string
	Username  string
	Database  string
}

func (secret ClientSecret) Key() string {
	return fmt.Sprintf("%s_%s", secret.Namespace, secret.Name)
}

func fromSecret(secret *corev1.Secret) (ClientSecret, error) {
	if secret == nil {
		return ClientSecret{}, errors.New("cannot parse Service from nil")
	}

	return ClientSecret{
		Namespace: secret.Namespace,
		Name:      secret.Name,
		Username:  string(secret.Data["username"]),
		Database:  string(secret.Data["database"]),
	}, nil
}

type ClientSecretAddedEvent struct {
	New ClientSecret
}

type ClientSecretUpdatedEvent struct {
	Old ClientSecret
	New ClientSecret
}

type ClientSecretDeletedEvent struct {
	Old ClientSecret
}

func (c *DeploymentsClient) ListenClientSecret(events chan<- interface{}) (cache.Store, chan<- struct{}) {
	serviceStore, serviceController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo v1.ListOptions) (result runtime.Object, err error) {
				lo.LabelSelector = "db-operator.ponglehub.co.uk/owned=true"
				return c.clientset.CoreV1().Secrets("").List(context.TODO(), lo)
			},
			WatchFunc: func(lo v1.ListOptions) (watch.Interface, error) {
				lo.LabelSelector = "db-operator.ponglehub.co.uk/owned=true"
				return c.clientset.CoreV1().Secrets("").Watch(context.TODO(), lo)
			},
		},
		&corev1.Secret{},
		time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(newObj interface{}) {
				newClientSecret := newObj.(*corev1.Secret)
				newClientSecretObj, err := fromSecret(newClientSecret)
				if err != nil {
					logrus.Errorf("Failed to parse new service: %+v", err)
					return
				}

				events <- ClientSecretAddedEvent{New: newClientSecretObj}
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				oldClientSecret := oldObj.(*corev1.Secret)
				oldClientSecretObj, err := fromSecret(oldClientSecret)
				if err != nil {
					logrus.Errorf("Failed to parse updated old service: %+v", err)
					return
				}

				newClientSecret := newObj.(*corev1.Secret)
				newClientSecretObj, err := fromSecret(newClientSecret)
				if err != nil {
					logrus.Errorf("Failed to parse updated new service: %+v", err)
					return
				}

				events <- ClientSecretUpdatedEvent{New: newClientSecretObj, Old: oldClientSecretObj}
			},
			DeleteFunc: func(obj interface{}) {
				oldClientSecret := obj.(*corev1.Secret)
				oldClientSecretObj, err := fromSecret(oldClientSecret)
				if err != nil {
					logrus.Errorf("Failed to parse deleted service: %+v", err)
					return
				}

				events <- ClientSecretDeletedEvent{Old: oldClientSecretObj}
			},
		},
	)

	stopper := make(chan struct{})
	go serviceController.Run(stopper)

	return serviceStore, stopper
}
