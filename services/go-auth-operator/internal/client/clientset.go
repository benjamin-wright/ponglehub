package client

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type UserClient struct {
	restClient rest.Interface
}

func New() (*UserClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kube config: %+v", err)
	}

	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create CRD rest client: %+v", err)
	}

	return &UserClient{restClient: client}, nil
}

func (c *UserClient) List(opts metav1.ListOptions) (*AuthUserList, error) {
	result := AuthUserList{}
	err := c.restClient.
		Get().
		Resource("authusers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *UserClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("authusers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}

func (c *UserClient) Listen(
	addFunc func(user *AuthUser),
	updateFunc func(oldUser *AuthUser, newUser *AuthUser),
	deleteFunc func(user *AuthUser),
) (cache.Store, chan<- struct{}) {
	userStore, userController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return c.List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return c.Watch(lo)
			},
		},
		&AuthUser{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				addFunc(obj.(*AuthUser))
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				updateFunc(oldObj.(*AuthUser), newObj.(*AuthUser))
			},
			DeleteFunc: func(obj interface{}) {
				deleteFunc(obj.(*AuthUser))
			},
		},
	)

	stopper := make(chan struct{})
	go userController.Run(stopper)

	return userStore, stopper
}

func (c *UserClient) Create(user AuthUser, opts metav1.CreateOptions) error {
	err := c.restClient.
		Post().
		Resource("authusers").
		Body(&user).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO())

	return err.Error()
}
