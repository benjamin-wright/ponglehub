package crds

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	restClient rest.Interface
}

type ClientArgs struct {
	External bool
}

func New(args *ClientArgs) (*Client, error) {
	var err error
	var config *rest.Config

	if args.External {
		KUBECONFIG, ok := os.LookupEnv("KUBECONFIG")
		if !ok {
			return nil, fmt.Errorf("failed to get kube config: missing KUBECONFIG env var")
		}

		config, err = clientcmd.BuildConfigFromFlags("", KUBECONFIG)
		if err != nil {
			return nil, fmt.Errorf("failed to get kube config: %+v", err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get kube config: %+v", err)
		}
	}

	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create CRD rest client: %+v", err)
	}

	return &Client{restClient: client}, nil
}

func (c *Client) list(opts v1.ListOptions) (*EventTriggerList, error) {
	result := EventTriggerList{}
	err := c.restClient.
		Get().
		Resource("eventtriggers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *Client) watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("eventtriggers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}

type ChangedHandler func(oldTrigger *EventTrigger, newTrigger *EventTrigger)

func (c *Client) Listen(handler ChangedHandler) (cache.Store, chan<- struct{}) {
	clientStore, clientController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo v1.ListOptions) (result runtime.Object, err error) {
				return c.list(lo)
			},
			WatchFunc: func(lo v1.ListOptions) (watch.Interface, error) {
				return c.watch(lo)
			},
		},
		&EventTrigger{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				trigger := obj.(*EventTrigger)
				handler(nil, trigger)
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				oldTrigger := oldObj.(*EventTrigger)
				newTrigger := newObj.(*EventTrigger)

				if oldTrigger.ResourceVersion != newTrigger.ResourceVersion {
					handler(oldTrigger, newTrigger)
				}
			},
			DeleteFunc: func(obj interface{}) {
				trigger := obj.(*EventTrigger)
				handler(trigger, nil)
			},
		},
	)

	stopper := make(chan struct{})
	go clientController.Run(stopper)

	return clientStore, stopper
}

func (c *Client) Create(name string, namespace string, filters []string, url string) error {
	trigger := EventTrigger{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: EventTriggerSpec{
			Filters: filters,
			URL:     url,
		},
	}

	err := c.restClient.
		Post().
		Namespace(namespace).
		Resource("eventtriggers").
		VersionedParams(&v1.CreateOptions{}, scheme.ParameterCodec).
		Body(&trigger).
		Do(context.TODO()).
		Error()

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Delete(name string, namespace string) error {
	return c.restClient.
		Delete().
		Namespace(namespace).
		Resource("eventtriggers").
		VersionedParams(&v1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).
		Error()
}
