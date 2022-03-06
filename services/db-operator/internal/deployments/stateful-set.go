package deployments

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type StatefulSet struct {
	Name      string
	Namespace string
	Storage   string
	Ready     bool
}

func fromSS(ss *appsv1.StatefulSet) (StatefulSet, error) {
	volumes := ss.Spec.VolumeClaimTemplates

	if len(volumes) != 1 {
		return StatefulSet{}, fmt.Errorf("bad database statefulset %s (%s), expected 1 volume got %d", ss.Name, ss.Namespace, len(volumes))
	}

	request := volumes[0].Spec.Resources.Requests.Storage()
	if request == nil {
		return StatefulSet{}, fmt.Errorf("bad database statefulset %s (%s), expected a storage request, got none", ss.Name, ss.Namespace)
	}

	return StatefulSet{
		Namespace: ss.Namespace,
		Name:      ss.Name,
		Storage:   request.String(),
		Ready:     ss.Status.Replicas == ss.Status.ReadyReplicas,
	}, nil
}

func (d *DeploymentsClient) ListStatefulSets() ([]StatefulSet, error) {
	statefulSetList, err := d.clientset.AppsV1().
		StatefulSets("").
		List(context.TODO(), v1.ListOptions{
			LabelSelector: "db-operator.ponglehub.co.uk/owned=true",
		})

	if err != nil {
		return nil, fmt.Errorf("failed to list database statefulsets: %+v", err)
	}

	statefulSets := []StatefulSet{}

	for _, statefulSet := range statefulSetList.Items {
		db, err := fromSS(&statefulSet)
		if err != nil {
			return nil, fmt.Errorf("failed to parse stateful set: %+v", err)
		}

		statefulSets = append(statefulSets, db)
	}

	return statefulSets, nil
}

func (d *DeploymentsClient) DeleteStatefulSet(statefulSet StatefulSet) error {
	logrus.Infof("Deleting stateful set %s (%s)", statefulSet.Name, statefulSet.Namespace)
	err := d.clientset.AppsV1().
		StatefulSets(statefulSet.Namespace).
		Delete(context.TODO(), statefulSet.Name, v1.DeleteOptions{})

	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete statefulset deployment: %+v", err)
	}

	pvcName := fmt.Sprintf("%s-%s-0", statefulSet.Name, statefulSet.Name)
	err = d.clientset.CoreV1().PersistentVolumeClaims(statefulSet.Namespace).Delete(context.TODO(), pvcName, v1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete statefulset PVC: %+v", err)
	}

	return nil
}

func (d *DeploymentsClient) AddStatefulSet(statefulSet StatefulSet) error {
	logrus.Infof("Creating stateful set %s (%s)", statefulSet.Name, statefulSet.Namespace)
	size, err := resource.ParseQuantity(statefulSet.Storage)
	if err != nil {
		return fmt.Errorf("failed to parse statefulset storage requirement: %+v", err)
	}

	statefulSetObject := appsv1.StatefulSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      statefulSet.Name,
			Namespace: statefulSet.Namespace,
			Labels: map[string]string{
				"db-operator.ponglehub.co.uk/owned": "true",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"db-operator.ponglehub.co.uk/deployment": statefulSet.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"db-operator.ponglehub.co.uk/deployment": statefulSet.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "db",
							Image:   "cockroachdb/cockroach:v20.2.8",
							Command: []string{"cockroach"},
							Args: []string{
								"--logtostderr",
								"start-single-node",
								"--insecure",
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 26257},
								{ContainerPort: 8080},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      statefulSet.Name,
									MountPath: "/cockroach/cockroach-data",
								},
							},
							ReadinessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/health?ready=1",
										Port: intstr.FromInt(8080),
									},
								},
								InitialDelaySeconds: 10,
								PeriodSeconds:       5,
								FailureThreshold:    2,
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: v1.ObjectMeta{
						Name:      statefulSet.Name,
						Namespace: statefulSet.Namespace,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								"storage": size,
							},
						},
					},
				},
			},
		},
	}

	_, err = d.clientset.AppsV1().StatefulSets(statefulSet.Namespace).Create(context.TODO(), &statefulSetObject, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to deploy statefulSet: %+v", err)
	}

	return nil
}

type StatefulSetAddedEvent struct {
	New StatefulSet
}

type StatefulSetUpdatedEvent struct {
	Old StatefulSet
	New StatefulSet
}

type StatefulSetDeletedEvent struct {
	Old StatefulSet
}

func (c *DeploymentsClient) ListenStatefulSets(events chan<- interface{}) (StatefulSetStore, chan<- struct{}) {
	ssStore, ssController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo v1.ListOptions) (result runtime.Object, err error) {
				lo.LabelSelector = "db-operator.ponglehub.co.uk/owned=true"
				return c.clientset.AppsV1().StatefulSets("").List(context.TODO(), lo)
			},
			WatchFunc: func(lo v1.ListOptions) (watch.Interface, error) {
				lo.LabelSelector = "db-operator.ponglehub.co.uk/owned=true"
				return c.clientset.AppsV1().StatefulSets("").Watch(context.TODO(), lo)
			},
		},
		&appsv1.StatefulSet{},
		time.Minute,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    func(newObj interface{}) { events <- true },
			UpdateFunc: func(oldObj interface{}, newObj interface{}) { events <- true },
			DeleteFunc: func(obj interface{}) { events <- true },
		},
	)

	stopper := make(chan struct{})
	go ssController.Run(stopper)

	return StatefulSetStore{store: ssStore}, stopper
}
