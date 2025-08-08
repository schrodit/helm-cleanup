package main

import (
	"github.com/spf13/pflag"
)

// EnvSettings defined settings
type EnvSettings struct {
	DryRun         bool
	Debug          bool
	KubeConfigFile string
	KubeContext    string
	Namespace      string
	Yes            bool
}

// New returns default env settings
func New() *EnvSettings {
	envSettings := EnvSettings{}
	return &envSettings
}

// AddFlags binds flags to the given flagset.
func (s *EnvSettings) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.DryRun, "dry-run", true, "simulate a command")
	fs.BoolVar(&s.Yes, "yes", false, "Do not prompt on every deleted resource")
	fs.StringVar(&s.KubeConfigFile, "kubeconfig", "", "path to the kubeconfig file")
	fs.StringVar(&s.KubeContext, "kube-context", s.KubeContext, "name of the kubeconfig context to use")
	fs.StringVar(&s.Namespace, "namespace", s.Namespace, "namespace scope of the release")
}
