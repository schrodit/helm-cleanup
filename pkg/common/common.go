package common

import (
	"fmt"
)

// KubeConfig are the Kubernetes configuration settings
type KubeConfig struct {
	Context string
	File    string
}

type Options struct {
	KubeConfig KubeConfig
	Namespace  string
	DryRun     bool
	Debug bool
}

type Release struct {
	Name      string
	Namespace string
}

func (r *Release) Key() string {
	return fmt.Sprintf("%s-%s", r.Name, r.Namespace)
}
