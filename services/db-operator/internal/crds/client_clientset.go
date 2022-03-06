package crds

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
)

type Client struct {
	Name       string
	Namespace  string
	Deployment string
	Database   string
	Username   string
	Secret     string
	Ready      bool
}

func (client Client) Key() string {
	return fmt.Sprintf("%s_%s", client.Namespace, client.Name)
}

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

func (c *DBClient) ClientList(namespace string) ([]Client, error) {
	result := CockroachClientList{}
	err := c.restClient.
		Get().
		Resource("cockroachclients").
		Namespace(namespace).
		VersionedParams(&v1.ListOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	clients := []Client{}
	for _, client := range result.Items {
		c, err := clientFromApi(&client)
		if err != nil {
			return nil, fmt.Errorf("error parsing client: %+v", err)
		}

		clients = append(clients, c)
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

func (c *DBClient) ClientCreate(client Client) error {
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

func (c *DBClient) ClientGet(name string, namespace string) (Client, error) {
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
		return Client{}, fmt.Errorf("error fetching cockroachclient: %+v", err)
	}

	clientObj, err := clientFromApi(&client)
	if err != nil {
		return Client{}, fmt.Errorf("error parsing cockroach client: %+v", err)
	}

	return clientObj, nil
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

func clientFromApi(client *CockroachClient) (Client, error) {
	if client == nil {
		return Client{}, errors.New("cannot parse client from nil")
	}

	return Client{
		Name:       client.Name,
		Username:   client.Spec.Username,
		Namespace:  client.Namespace,
		Deployment: client.Spec.Deployment,
		Database:   client.Spec.Database,
		Ready:      client.Status.Ready,
	}, nil
}

func (c *DBClient) ClientListen(events chan<- interface{}) (cache.Store, chan<- struct{}) {
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
		time.Second*20,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(newObj interface{}) { events <- true },
			UpdateFunc: func(oldObj interface{}, newObj interface{}) { events <- true },
			DeleteFunc: func(oldObj interface{}) { events <- true },
		},
	)

	stopper := make(chan struct{})
	go clientController.Run(stopper)

	return clientStore, stopper
}
