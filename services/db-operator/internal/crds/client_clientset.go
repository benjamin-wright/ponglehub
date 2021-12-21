package crds

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
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

func (c *DBClient) ClientList(namespace string) ([]types.Client, error) {
	result := CockroachClientList{}
	err := c.restClient.
		Get().
		Resource("cockroachclients").
		Namespace(namespace).
		VersionedParams(&v1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	clients := []types.Client{}
	for _, client := range result.Items {
		clients = append(clients, clientFromApi(&client))
	}

	return clients, err
}

func (c *DBClient) clientWatch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("cockroachclients").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}

func (c *DBClient) ClientCreate(client types.Client) error {
	clientObj := CockroachClient{
		ObjectMeta: v1.ObjectMeta{
			Name:      client.Name,
			Namespace: client.Namespace,
		},
		Spec: CockroachClientSpec{
			Deployment: client.Deployment,
			Database:   client.Database,
			Username:   client.Username,
		},
		Status: CockroachClientStatus{
			Ready: client.Ready,
		},
	}

	err := c.restClient.
		Post().
		Namespace(client.Namespace).
		Resource("cockroachclients").
		VersionedParams(&v1.CreateOptions{}, scheme.ParameterCodec).
		Body(&clientObj).
		Do(context.TODO()).
		Error()

	if err != nil {
		return err
	}

	return nil
}

func (c *DBClient) ClientGet(name string, namespace string) (types.Client, error) {
	client := CockroachClient{}

	err := c.restClient.
		Get().
		Namespace(namespace).
		Resource("cockroachclients").
		Name(name).
		VersionedParams(&v1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&client)

	if err != nil {
		return types.Client{}, err
	}

	return clientFromApi(&client), nil
}

func (c *DBClient) ClientDelete(name string, namespace string) error {
	return c.restClient.
		Delete().
		Namespace(namespace).
		Resource("cockroachclients").
		Name(name).
		VersionedParams(&v1.DeleteOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Error()
}

func (c *DBClient) ClientUpdate(name string, namespace string, ready bool) error {
	client := CockroachClient{}

	err := c.restClient.
		Get().
		Namespace(namespace).
		Resource("cockroachclients").
		Name(name).
		VersionedParams(&v1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&client)

	if err != nil {
		return fmt.Errorf("failed in initial fetch: %+v", err)
	}

	if client.Status.Ready == ready {
		return nil
	}

	logrus.Infof("Updating Client CRD status: %s (%s) -> ready: %t", name, namespace, ready)
	client.Status.Ready = ready

	return c.restClient.
		Put().
		Namespace(namespace).
		Resource("cockroachclients").
		Name(name).
		SubResource("status").
		VersionedParams(&v1.UpdateOptions{}, scheme.ParameterCodec).
		Body(&client).
		Do(context.TODO()).
		Error()
}

func clientFromApi(client *CockroachClient) types.Client {
	return types.Client{
		Name:       client.Name,
		Username:   client.Spec.Username,
		Namespace:  client.Namespace,
		Deployment: client.Spec.Deployment,
		Database:   client.Spec.Database,
		Ready:      client.Status.Ready,
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
