package users

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

type User struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ID         string `json:"id"`
	Pending    bool   `json:"pending"`
	Generation int64  `json:"generation"`
}

func (a User) Equals(user User) bool {
	return a.Email == user.Email &&
		a.Username == user.Username &&
		a.Password == user.Password
}

func fromUser(user User) *AuthUser {
	return &AuthUser{
		ObjectMeta: v1.ObjectMeta{
			Name:       user.Name,
			Generation: user.Generation,
		},
		Spec: AuthUserSpec{
			Name:     user.Username,
			Email:    user.Email,
			Password: user.Password,
		},
		Status: AuthUserStatus{
			ID:      user.ID,
			Pending: user.Pending,
		},
	}
}

func fromAuthUser(authUser *AuthUser) User {
	return User{
		Name:       authUser.Name,
		Username:   authUser.Spec.Name,
		Email:      authUser.Spec.Email,
		Password:   authUser.Spec.Password,
		ID:         authUser.Status.ID,
		Pending:    authUser.Status.Pending,
		Generation: authUser.Generation,
	}
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

func (c *UserClient) Create(user User) error {
	authUser := fromUser(user)
	res := c.restClient.
		Post().
		Resource("authusers").
		VersionedParams(&v1.CreateOptions{}, scheme.ParameterCodec).
		Body(authUser).
		Do(context.TODO())

	return res.Error()
}

func (c *UserClient) Update(user User) error {
	authUser := fromUser(user)
	return c.restClient.
		Put().
		Resource("authusers").
		Name(user.Name).
		VersionedParams(&v1.UpdateOptions{}, scheme.ParameterCodec).
		Body(authUser).
		Do(context.TODO()).
		Error()
}

func (c *UserClient) Status(user User) error {
	authUser := fromUser(user)
	return c.restClient.
		Put().
		Resource("authusers").
		Name(user.Name).
		SubResource("status").
		VersionedParams(&v1.UpdateOptions{}, scheme.ParameterCodec).
		Body(authUser).
		Do(context.TODO()).
		Error()
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
				updateFunc(
					fromAuthUser(oldObj.(*AuthUser)),
					fromAuthUser(newObj.(*AuthUser)),
				)
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
