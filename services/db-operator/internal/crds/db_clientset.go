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

func (c *DBClient) dbList(opts v1.ListOptions) (*CockroachDBList, error) {
	result := CockroachDBList{}
	err := c.restClient.
		Get().
		Resource("cockroachdbs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *DBClient) dbWatch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Resource("cockroachdbs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}

func dbFromApi(db *CockroachDB) types.Database {
	return types.Database{
		Name:      db.Name,
		Namespace: db.Namespace,
		Storage:   db.Spec.Storage,
	}
}

func apiFromDB(db types.Database) *CockroachDB {
	return &CockroachDB{
		ObjectMeta: v1.ObjectMeta{
			Name:      db.Name,
			Namespace: db.Namespace,
		},
		Spec: CockroachDBSpec{
			Storage: db.Storage,
		},
	}
}

type DBAddedHandler func(client types.Database)
type DBUpdatedHandler func(oldDB types.Database, newDB types.Database)
type DBDeletedHandler func(client types.Database)

func (c *DBClient) DBListen(added DBAddedHandler, updated DBUpdatedHandler, deleted DBDeletedHandler) (cache.Store, chan<- struct{}) {
	dbStore, dbController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo v1.ListOptions) (result runtime.Object, err error) {
				return c.dbList(lo)
			},
			WatchFunc: func(lo v1.ListOptions) (watch.Interface, error) {
				return c.dbWatch(lo)
			},
		},
		&CockroachDB{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				added(dbFromApi(obj.(*CockroachDB)))
			},
			UpdateFunc: func(oldObj interface{}, newObj interface{}) {
				oldClient := oldObj.(*CockroachDB)
				newClient := newObj.(*CockroachDB)

				if oldClient.Generation != newClient.Generation {
					updated(dbFromApi(oldClient), dbFromApi(newClient))
				}
			},
			DeleteFunc: func(obj interface{}) {
				deleted(dbFromApi(obj.(*CockroachDB)))
			},
		},
	)

	stopper := make(chan struct{})
	go dbController.Run(stopper)

	return dbStore, stopper
}

func (c *DBClient) DBCreate(db types.Database) error {
	dbObj := apiFromDB(db)

	return c.restClient.
		Post().
		Namespace(db.Namespace).
		Resource("cockroachdbs").
		VersionedParams(&v1.CreateOptions{}, scheme.ParameterCodec).
		Body(dbObj).
		Do(context.TODO()).
		Error()
}

func (c *DBClient) DBGet(name string, namespace string) (types.Database, error) {
	var db CockroachDB

	err := c.restClient.
		Get().
		Namespace(namespace).
		Resource("cockroachdbs").
		VersionedParams(&v1.CreateOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).
		Into(&db)

	if err != nil {
		return types.Database{}, err
	}

	return dbFromApi(&db), nil
}

func (c *DBClient) DBDelete(name string, namespace string) error {
	return c.restClient.
		Delete().
		Namespace(namespace).
		Resource("cockroachdbs").
		VersionedParams(&v1.CreateOptions{}, scheme.ParameterCodec).
		Name(name).
		Do(context.TODO()).
		Error()
}
