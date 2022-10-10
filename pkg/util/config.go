package util

import (
	"errors"
	"path/filepath"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	FilePath   string
	Namespace  string
	KubeConfig string
)

// NewKubeConfigClient returns the kubeconfig and the client created from the kubeconfig.
func NewKubeConfigClient() (*rest.Config, *k8s.Clientset, error) {
	kubeconfig := ""
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	if KubeConfig != "" {
		kubeconfig = KubeConfig
	}
	if kubeconfig == "" {
		return nil, nil, errors.New("kubeconfig is required")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	client, err := k8s.NewForConfig(config)
	if err != nil {
		return config, nil, err
	}
	return config, client, nil
}
