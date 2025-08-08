package k8s

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/schrodit/helm-cleanup/pkg/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func DeleteUnstrcutured(ctx context.Context, kc *common.KubeClient, r *common.KubeResource) error {
	resource := kc.Dynamic.Resource(r.GroupVersionResource)
	if r.GetNamespace() != "" {
		return resource.Namespace(r.GetNamespace()).Delete(ctx, r.GetName(), metav1.DeleteOptions{})
	}
	return resource.Delete(ctx, r.GetName(), metav1.DeleteOptions{})
}

// GetClientSetWithKubeConfig returns a kubernetes ClientSet
func GetClientSetWithKubeConfig(kubeConfigFile, context string) (*common.KubeClient, error) {
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
		return nil, errors.Wrap(err, "failed to build config from flags")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubeclient")
	}
	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dynamic kubeclient")
	}

	return &common.KubeClient{
		Default: clientset,
		Dynamic: dyn,
	}, nil
}

func buildConfigFromFlags(context string, kubeConfigFiles []string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{Precedence: kubeConfigFiles},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}
