package common

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// KubeConfig are the Kubernetes configuration settings
type KubeConfig struct {
	Context string
	File    string
}

type KubeClient struct {
	Default *kubernetes.Clientset
	Dynamic *dynamic.DynamicClient
}

type KubeResource struct {
	*unstructured.Unstructured
	GroupVersionResource schema.GroupVersionResource
}

type Options struct {
	KubeConfig KubeConfig
	Namespace  string
	DryRun     bool
	Debug      bool
}

type Release struct {
	Name      string
	Namespace string
}

func (r *Release) Key() string {
	return fmt.Sprintf("%s/%s", r.Name, r.Namespace)
}
