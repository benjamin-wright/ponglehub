package deployments

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodeCerts struct {
	CACrt   []byte
	NodeCrt []byte
	NodeKey []byte
}

var emptyCert = NodeCerts{}

func (d *DeploymentsClient) GetNodeSecret(namespace string, name string) (NodeCerts, error) {
	secret, err := d.clientset.CoreV1().Secrets(namespace).Get(context.TODO(), name+"-ssl", v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return emptyCert, nil
		}

		return emptyCert, fmt.Errorf("failed to get Node secret: %+v", err)
	}

	if err := mapHas(secret.Data, []string{"ca.crt", "node.crt", "node.key"}); err != nil {
		return emptyCert, fmt.Errorf("failed to parse secret: %+v", err)
	}

	return NodeCerts{
		CACrt:   secret.Data["ca.crt"],
		NodeCrt: secret.Data["node.crt"],
		NodeKey: secret.Data["node.key"],
	}, nil
}

func (d *DeploymentsClient) DeleteNodeSecret(namespace string, name string) error {
	err := d.clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), name+"-ssl", v1.DeleteOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}

		return fmt.Errorf("failed to delete node secret: %+v", err)
	}

	return nil
}

func (d *DeploymentsClient) AddNodeSecret(namespace string, name string, certs NodeCerts) error {
	secret := corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      name + "-ssl",
			Namespace: namespace,
			Labels: map[string]string{
				"db-operator.ponglehub.co.uk/owned": "true",
			},
		},
		Data: map[string][]byte{
			"ca.crt":   certs.CACrt,
			"node.crt": certs.NodeCrt,
			"node.key": certs.NodeKey,
		},
	}

	_, err := d.clientset.CoreV1().Secrets(namespace).Create(context.TODO(), &secret, v1.CreateOptions{})

	return err
}
