package crds

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type UserClient struct {
	restClient rest.Interface
}

func fromUser(user User) *AuthUser {
	authUser := AuthUser{
		ObjectMeta: v1.ObjectMeta{
			Name:            user.Name,
			ResourceVersion: user.ResourceVersion,
		},
		Spec: AuthUserSpec{
			Display:      user.Display,
			Email:        user.Email,
			InviteToken:  user.InviteToken,
			PasswordHash: user.PasswordHash,
		},
	}

	if user.ID != "" {
		authUser.ObjectMeta.UID = types.UID(user.ID)
	}

	return &authUser
}

func fromAuthUser(authUser *AuthUser) User {
	return User{
		ID:              string(authUser.UID),
		Name:            authUser.Name,
		ResourceVersion: authUser.ResourceVersion,
		Display:         authUser.Spec.Display,
		Email:           authUser.Spec.Email,
		InviteToken:     authUser.Spec.InviteToken,
		PasswordHash:    authUser.Spec.PasswordHash,
	}
}

type ClientArgs struct {
	External bool
}

func New(args *ClientArgs) (*UserClient, error) {
	var config *rest.Config
	var err error

	if args != nil && args.External {
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

	return &UserClient{restClient: client}, nil
}

func (c *UserClient) Get(name string) (User, error) {
	result := AuthUser{}
	err := c.restClient.
		Get().
		Resource("authusers").
		Name(name).
		VersionedParams(&v1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return fromAuthUser(&result), err
}

func (c *UserClient) Delete(name string) error {
	res := c.restClient.
		Delete().
		Resource("authusers").
		VersionedParams(&v1.DeleteOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO())

	return res.Error()
}

func (c *UserClient) Create(user User) (User, error) {
	authUser := fromUser(user)

	result := AuthUser{}
	err := c.restClient.
		Post().
		Resource("authusers").
		VersionedParams(&v1.CreateOptions{}, scheme.ParameterCodec).
		Body(authUser).
		Do(context.TODO()).
		Into(&result)

	return fromAuthUser(&result), err
}

func (c *UserClient) Update(user User) (User, error) {
	authUser := fromUser(user)

	result := AuthUser{}
	err := c.restClient.
		Put().
		Resource("authusers").
		Name(user.Name).
		VersionedParams(&v1.UpdateOptions{}, scheme.ParameterCodec).
		Body(authUser).
		Do(context.TODO()).
		Into(&result)

	return fromAuthUser(&result), err
}

func (c *UserClient) list(opts v1.ListOptions) (*AuthUserList, error) {
	result := AuthUserList{}
	err := c.restClient.
		Get().
		Resource("authusers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *UserClient) watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("authusers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}

func (c *UserClient) Listen(
	addFunc func(user User),
	updateFunc func(oldUser User, newUser User),
	deleteFunc func(user User),
) (cache.Store, chan<- struct{}) {
	userStore, userController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo v1.ListOptions) (result runtime.Object, err error) {
				return c.list(lo)
			},
			WatchFunc: func(lo v1.ListOptions) (watch.Interface, error) {
				return c.watch(lo)
			},
		},
		&AuthUser{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				addFunc(
					fromAuthUser(obj.(*AuthUser)),
				)
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				oldUser := oldObj.(*AuthUser)
				newUser := newObj.(*AuthUser)

				if oldUser.Generation != newUser.Generation {
					updateFunc(
						fromAuthUser(oldUser),
						fromAuthUser(newUser),
					)
				}
			},
			DeleteFunc: func(obj interface{}) {
				deleteFunc(
					fromAuthUser(obj.(*AuthUser)),
				)
			},
		},
	)

	stopper := make(chan struct{})
	go userController.Run(stopper)

	return userStore, stopper
}
