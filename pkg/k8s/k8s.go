package k8s

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetClientSetWithKubeConfig returns a kubernetes ClientSet
func GetClientSetWithKubeConfig(kubeConfigFile, context string) (*kubernetes.Clientset, *dynamic.DynamicClient, error) {
	var kubeConfigFiles []string
	if kubeConfigFile != "" {
		kubeConfigFiles = append(kubeConfigFiles, kubeConfigFile)
	} else if kubeConfigPath := os.Getenv("KUBECONFIG"); kubeConfigPath != "" {
		// The KUBECONFIG environment variable holds a list of kubeconfig files.
		// For Linux and Mac, the list is colon-delimited. For Windows, the list
		// is semicolon-delimited. Ref:
		// https://kubernetes.io/docs/concepts/configuration/organize-cluster-access-kubeconfig/#the-kubeconfig-environment-variable
		var separator string
		if runtime.GOOS == "windows" {
			separator = ";"
		} else {
			separator = ":"
		}
		kubeConfigFiles = strings.Split(kubeConfigPath, separator)
	} else {
		kubeConfigFiles = append(kubeConfigFiles, filepath.Join(os.Getenv("HOME"), ".kube", "config"))
	}

	config, err := buildConfigFromFlags(context, kubeConfigFiles)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to build config from flags")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create kubeclient")
	}
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create dynamic kubeclient")
	}

	return clientset, dyn, nil
}

func buildConfigFromFlags(context string, kubeConfigFiles []string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{Precedence: kubeConfigFiles},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}
