package deployments

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"ponglehub.co.uk/operators/db/internal/types"
)

var emptyDB = types.Database{}

func fromSS(ss appsv1.StatefulSet) (types.Database, error) {
	volumes := ss.Spec.VolumeClaimTemplates

	if len(volumes) != 1 {
		return emptyDB, fmt.Errorf("bad database deployment %s (%s), expected 1 volume got %d", ss.Name, ss.Namespace, len(volumes))
	}

	request := volumes[0].Spec.Resources.Requests.Storage()
	if request == nil {
		return emptyDB, fmt.Errorf("bad database deployment %s (%s), expected a storage request, got none", ss.Name, ss.Namespace)
	}

	return types.Database{
		Namespace: ss.Namespace,
		Name:      ss.Name,
		Storage:   request.String(),
	}, nil
}

func (d *DeploymentsClient) GetDeployments(namespace string) ([]types.Database, error) {
	deployments, err := d.clientset.AppsV1().
		StatefulSets(namespace).
		List(context.TODO(), v1.ListOptions{
			LabelSelector: "db-operator.ponglehub.co.uk/owned=true",
		})

	if err != nil {
		return nil, fmt.Errorf("failed to list database deployments: %+v", err)
	}

	databases := []types.Database{}

	for _, deployment := range deployments.Items {
		db, err := fromSS(deployment)
		if err != nil {
			return nil, fmt.Errorf("failed to parse stateful set: %+v", err)
		}

		databases = append(databases, db)
	}

	return databases, nil
}

func (d *DeploymentsClient) DeleteDeployment(database types.Database) error {
	err := d.clientset.AppsV1().
		StatefulSets(database.Namespace).
		Delete(context.TODO(), database.Name, v1.DeleteOptions{})

	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete database deployment: %+v", err)
	}

	pvcName := fmt.Sprintf("%s-%s-0", database.Name, database.Name)
	err = d.clientset.CoreV1().PersistentVolumeClaims(database.Namespace).Delete(context.TODO(), pvcName, v1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("failed to delete database PVC: %+v", err)
	}

	return nil
}

func (d *DeploymentsClient) AddDeployment(database types.Database) error {
	size, err := resource.ParseQuantity(database.Storage)
	if err != nil {
		return fmt.Errorf("failed to parse database storage requirement: %+v", err)
	}

	var certReadMode int32 = 256

	deployment := appsv1.StatefulSet{
		ObjectMeta: v1.ObjectMeta{
			Name:      database.Name,
			Namespace: database.Namespace,
			Labels: map[string]string{
				"db-operator.ponglehub.co.uk/owned": "true",
			},
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &v1.LabelSelector{
				MatchLabels: map[string]string{
					"db-operator.ponglehub.co.uk/deployment": database.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"db-operator.ponglehub.co.uk/deployment": database.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "db",
							Image:   "cockroachdb/cockroach:v20.2.8",
							Command: []string{"cockroach"},
							Args:    []string{"start-single-node", "--certs-dir", "/certs"},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 26257},
								{ContainerPort: 8080},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      database.Name,
									MountPath: "/cockroach/cockroach-data",
								},
								{
									Name:      "ssl",
									MountPath: "/certs",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "ssl",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName:  database.Name + "-ssl",
									DefaultMode: &certReadMode,
								},
							},
						},
					},
				},
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: v1.ObjectMeta{
						Name:      database.Name,
						Namespace: database.Namespace,
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

	_, err = d.clientset.AppsV1().StatefulSets(database.Namespace).Create(context.TODO(), &deployment, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to deploy database: %+v", err)
	}

	return nil
}
