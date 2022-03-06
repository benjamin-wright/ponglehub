package crds

import (
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type DBClient struct {
	restClient rest.Interface
}

func New() (*DBClient, error) {
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

	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: GroupName, Version: GroupVersion}
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create CRD rest client: %+v", err)
	}

	return &DBClient{restClient: client}, nil
}
