package deployments

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type DeploymentsClient struct {
	clientset *kubernetes.Clientset
}

func New() (*DeploymentsClient, error) {
	var config *rest.Config
	var err error

	if KUBECONFIG, ok := os.LookupEnv("KUBECONFIG"); ok {
		config, err = clientcmd.BuildConfigFromFlags("", KUBECONFIG)
		if err != nil {
			return nil, fmt.Errorf("failed to load kube config from file: %+v", err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to get in-cluster kube config: %+v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube client: %+v", err)
	}

	return &DeploymentsClient{
		clientset: clientset,
	}, nil
}
