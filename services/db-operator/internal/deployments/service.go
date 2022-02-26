package deployments

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type Service struct {
	Namespace string
	Name      string
}

func (service Service) Key() string {
	return fmt.Sprintf("%s_%s", service.Namespace, service.Name)
}

func fromService(service *corev1.Service) (Service, error) {
	if service == nil {
		return Service{}, errors.New("cannot parse Service from nil")
	}

	return Service{
		Namespace: service.Namespace,
		Name:      service.Name,
	}, nil
}

func (d *DeploymentsClient) HasService(config Service) (bool, error) {
	_, err := d.clientset.CoreV1().Services(config.Namespace).Get(context.TODO(), config.Name, v1.GetOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (d *DeploymentsClient) AddService(config Service) error {
	service := corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name: config.Name,
			Labels: map[string]string{
				"db-operator.ponglehub.co.uk/owned": "true",
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"db-operator.ponglehub.co.uk/deployment": config.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "grpc",
					Port:       26257,
					TargetPort: intstr.FromInt(26257),
				},
				{
					Name:       "http",
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}

	_, err := d.clientset.CoreV1().Services(config.Namespace).Create(context.TODO(), &service, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service for cockroach deployment %s (%s): %+v", config.Name, config.Namespace, err)
	}

	return nil
}

func (d *DeploymentsClient) DeleteService(config Service) error {
	err := d.clientset.CoreV1().Services(config.Namespace).Delete(context.TODO(), config.Name, v1.DeleteOptions{})
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return nil
		}

		return err
	}

	return nil
}

type ServiceAddedEvent struct {
	New Service
}

type ServiceUpdatedEvent struct {
	Old Service
	New Service
}

type ServiceDeletedEvent struct {
	Old Service
}

func (c *DeploymentsClient) ListenService(events chan<- interface{}) (cache.Store, chan<- struct{}) {
	serviceStore, serviceController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo v1.ListOptions) (result runtime.Object, err error) {
				lo.LabelSelector = "db-operator.ponglehub.co.uk/owned=true"
				return c.clientset.CoreV1().Services("").List(context.TODO(), lo)
			},
			WatchFunc: func(lo v1.ListOptions) (watch.Interface, error) {
				lo.LabelSelector = "db-operator.ponglehub.co.uk/owned=true"
				return c.clientset.CoreV1().Services("").Watch(context.TODO(), lo)
			},
		},
		&corev1.Service{},
		time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(newObj interface{}) {
				newService := newObj.(*corev1.Service)
				newServiceObj, err := fromService(newService)
				if err != nil {
					logrus.Errorf("Failed to parse new service: %+v", err)
					return
				}

				events <- ServiceAddedEvent{New: newServiceObj}
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				oldService := oldObj.(*corev1.Service)
				oldServiceObj, err := fromService(oldService)
				if err != nil {
					logrus.Errorf("Failed to parse updated old service: %+v", err)
					return
				}

				newService := newObj.(*corev1.Service)
				newServiceObj, err := fromService(newService)
				if err != nil {
					logrus.Errorf("Failed to parse updated new service: %+v", err)
					return
				}

				events <- ServiceUpdatedEvent{New: newServiceObj, Old: oldServiceObj}
			},
			DeleteFunc: func(obj interface{}) {
				oldService := obj.(*corev1.Service)
				oldServiceObj, err := fromService(oldService)
				if err != nil {
					logrus.Errorf("Failed to parse deleted service: %+v", err)
					return
				}

				events <- ServiceDeletedEvent{Old: oldServiceObj}
			},
		},
	)

	stopper := make(chan struct{})
	go serviceController.Run(stopper)

	return serviceStore, stopper
}
