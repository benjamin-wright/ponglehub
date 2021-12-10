package crds

import (
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"ponglehub.co.uk/operators/db/internal/types"
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

func clientFromApi(old *CockroachClient) types.Client {
	return types.Client{
		Name:      old.Name,
		Namespace: old.Namespace,
		Database:  old.Spec.Database,
		Secret:    old.Spec.Secret,
	}
}

type ClientAddedHandler func(client types.Client)
type ClientUpdatedHandler func(oldClient types.Client, newClient types.Client)
type ClientDeletedHandler func(client types.Client)

func (c *DBClient) ClientListen(added ClientAddedHandler, updated ClientUpdatedHandler, deleted ClientDeletedHandler) (cache.Store, chan<- struct{}) {
	clientStore, clientController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo v1.ListOptions) (result runtime.Object, err error) {
				return c.clientList(lo)
			},
			WatchFunc: func(lo v1.ListOptions) (watch.Interface, error) {
				return c.clientWatch(lo)
			},
		},
		&CockroachClient{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				added(clientFromApi(obj.(*CockroachClient)))
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				oldClient := oldObj.(*CockroachClient)
				newClient := newObj.(*CockroachClient)

				if oldClient.Generation != newClient.Generation {
					updated(clientFromApi(oldClient), clientFromApi(newClient))
				}
			},
			DeleteFunc: func(obj interface{}) {
				deleted(clientFromApi(obj.(*CockroachClient)))
			},
		},
	)

	stopper := make(chan struct{})
	go clientController.Run(stopper)

	return clientStore, stopper
}
