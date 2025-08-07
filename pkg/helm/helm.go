package helm

import (
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"

	"github.com/pkg/errors"
	common "github.com/schrodit/helm-cleanup/pkg/common"
)

// Helm labels from https://github.com/helm/helm/blob/d7df660418a90bb93d87561efe2f503b60069c49/pkg/action/validate.go#L37
const (
	AppManagedByLabel              = "app.kubernetes.io/managed-by"
	AppManagedByHelm               = "Helm"
	HelmReleaseNameAnnotation      = "meta.helm.sh/release-name"
	HelmReleaseNamespaceAnnotation = "meta.helm.sh/release-namespace"
)

var (
	settings = cli.New()
)

func ListReleases(opts common.Options) ([]common.Release, error) {
	cfg, err := GetActionConfig("", opts.KubeConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get Helm action configuration")
	}

	releases, err := cfg.Releases.List(func(*release.Release) bool { return true })
	if err != nil {
		return nil, errors.Wrap(err, "failed to list helm releases")
	}

	res := make([]common.Release, len(releases))
	for i, r := range releases {
		res[i] = common.Release{
			Name:      r.Name,
			Namespace: r.Namespace,
		}
	}
	return res, nil
}

// GetActionConfig returns action configuration based on Helm env
func GetActionConfig(namespace string, kubeConfig common.KubeConfig) (*action.Configuration, error) {
	actionConfig := &action.Configuration{}

	// Add kube config settings passed by user
	settings.KubeConfig = kubeConfig.File
	settings.KubeContext = kubeConfig.Context

	err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debug)
	if err != nil {
		return nil, err
	}

	return actionConfig, err
}

func debug(format string, v ...interface{}) {
	if settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		err := log.Output(2, fmt.Sprintf(format, v...))
		if err != nil {
			return
		}
	}
}
