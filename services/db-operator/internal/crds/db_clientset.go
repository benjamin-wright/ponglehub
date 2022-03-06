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

type Database struct {
	Name      string
	Namespace string
	Storage   string
	Ready     bool
}

func (database Database) Key() string {
	return fmt.Sprintf("%s_%s", database.Namespace, database.Name)
}

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

func dbFromApi(db *CockroachDB) (Database, error) {
	if db == nil {
		return Database{}, errors.New("cannot parse database from nil")
	}

	return Database{
		Name:      db.Name,
		Namespace: db.Namespace,
		Storage:   db.Spec.Storage,
		Ready:     db.Status.Ready,
	}, nil
}

func apiFromDB(db Database) *CockroachDB {
	return &CockroachDB{
		ObjectMeta: v1.ObjectMeta{
			Name:      db.Name,
			Namespace: db.Namespace,
		},
		Spec: CockroachDBSpec{
			Storage: db.Storage,
		},
		Status: CockroachDBStatus{
			Ready: db.Ready,
		},
	}
}

func (c *DBClient) DBCreate(db Database) error {
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

func (c *DBClient) DBGet(name string, namespace string) (Database, error) {
	db := CockroachDB{}

	err := c.restClient.
		Get().
		Namespace(namespace).
		Resource("cockroachdbs").
		Name(name).
		VersionedParams(&v1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&db)

	if err != nil {
		return Database{}, err
	}

	dbObj, err := dbFromApi(&db)
	if err != nil {
		return Database{}, fmt.Errorf("error parsing database crd: %+v", err)
	}

	return dbObj, nil
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

func (c *DBClient) DBUpdate(name string, namespace string, ready bool) error {
	db := CockroachDB{}

	err := c.restClient.
		Get().
		Namespace(namespace).
		Resource("cockroachdbs").
		Name(name).
		VersionedParams(&v1.GetOptions{}, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&db)

	if err != nil {
		return fmt.Errorf("failed in initial fetch: %+v", err)
	}

	if db.Status.Ready == ready {
		return nil
	}

	logrus.Infof("Updating CRD status: %s (%s) -> ready: %t", name, namespace, ready)
	db.Status.Ready = ready

	return c.restClient.
		Put().
		Namespace(namespace).
		Resource("cockroachdbs").
		Name(name).
		SubResource("status").
		VersionedParams(&v1.UpdateOptions{}, scheme.ParameterCodec).
		Body(&db).
		Do(context.TODO()).
		Error()
}

func (c *DBClient) DBListen(events chan<- interface{}) (cache.Store, chan<- struct{}) {
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
		time.Second*10,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(newObj interface{}) { events <- true },
			UpdateFunc: func(oldObj interface{}, newObj interface{}) { events <- true },
			DeleteFunc: func(oldObj interface{}) { events <- true },
		},
	)

	stopper := make(chan struct{})
	go dbController.Run(stopper)

	return dbStore, stopper
}
