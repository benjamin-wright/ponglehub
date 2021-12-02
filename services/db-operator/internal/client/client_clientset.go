package client

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
)

func (c *DBClient) clientList(opts v1.ListOptions) (*CockroachClientList, error) {
	result := CockroachClientList{}
	err := c.restClient.
		Get().
		Resource("cockroachclients").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *DBClient) clientWatch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("cockroachclients").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}

func (c *DBClient) ClientListen() (cache.Store, chan<- struct{}) {
	userStore, userController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo v1.ListOptions) (result runtime.Object, err error) {
				return c.clientList(lo)
			},
			WatchFunc: func(lo v1.ListOptions) (watch.Interface, error) {
				return c.clientWatch(lo)
			},
		},
		&CockroachClient{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(obj interface{}) {},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {},
			DeleteFunc: func(obj interface{}) {},
		},
	)

	stopper := make(chan struct{})
	go userController.Run(stopper)

	return userStore, stopper
}
