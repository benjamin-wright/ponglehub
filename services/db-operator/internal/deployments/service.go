package deployments

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (d *DeploymentsClient) HasService(namespace string, name string) (bool, error) {
	_, err := d.clientset.CoreV1().Services(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (d *DeploymentsClient) AddService(namespace string, name string) error {
	service := corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"db-operator.ponglehub.co.uk/deployment": name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "client",
					Port:       26257,
					TargetPort: intstr.FromInt(26257),
				},
				{
					Name:       "web",
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}

	_, err := d.clientset.CoreV1().Services(namespace).Create(context.TODO(), &service, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service for cockroach deployment %s (%s): %+v", name, namespace, err)
	}

	return nil
}

func (d *DeploymentsClient) DeleteService(namespace string, name string) error {
	err := d.clientset.CoreV1().Services(namespace).Delete(context.TODO(), name, v1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return err
	}

	return nil
}
