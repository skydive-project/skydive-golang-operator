package config

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func GetConfig() (*rest.Config, error) {
	var config *rest.Config
	var err error
	if kubeConfig := os.Getenv("KUBECONFIG"); kubeConfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else if homePath := os.Getenv("HOME"); homePath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", homePath+"/.kube/config")
		if err != nil {
			config, err = rest.InClusterConfig()
		}
	} else {
		config, err = rest.InClusterConfig()
	}
	return config, err
}
