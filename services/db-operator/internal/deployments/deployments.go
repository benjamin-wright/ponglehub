package deployments

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type DeploymentsClient struct {
	clientset *kubernetes.Clientset
}

func New() (*DeploymentsClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kube config: %+v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube client: %+v", err)
	}

	return &DeploymentsClient{
		clientset: clientset,
	}, nil
}
