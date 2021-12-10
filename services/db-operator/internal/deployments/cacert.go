package deployments

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (d *DeploymentsClient) GetCASecret(namespace string, name string) ([]byte, error) {
	secret, err := d.clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name+"-ca", v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to get CA secret: %+v", err)
	}

	if key, ok := secret.Data["ca.key"]; ok {
		return key, nil
	} else {
		return nil, fmt.Errorf("failed to get CA secret: missing key ca.key")
	}
}

func (d *DeploymentsClient) DeleteCASecret(namespace string, name string) error {
	err := d.clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), name+"-ca", v1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return fmt.Errorf("failed to delete CA secret: %+v", err)
	}

	return nil
}

func (d *DeploymentsClient) AddCaSecret(namespace string, name string, key []byte) error {
	secret := corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      name + "-ca",
			Namespace: namespace,
			Labels: map[string]string{
				"db-operator.ponglehub.co.uk/owned": "true",
			},
		},
		Data: map[string][]byte{
			"ca.key": key,
		},
	}

	_, err := d.clientset.CoreV1().Secrets(namespace).Create(context.TODO(), &secret, v1.CreateOptions{})

	return err
}
