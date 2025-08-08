package cleanup

import (
	"context"

	"github.com/schrodit/helm-cleanup/pkg/common"
	"github.com/schrodit/helm-cleanup/pkg/helm"
	"k8s.io/apimachinery/pkg/util/sets"
)

func ListLeakedResources(
	ctx context.Context,
	releases []common.Release,
	kc *common.KubeClient,
	opts common.Options,
) ([]*common.KubeResource, error) {
	releasesSet := sets.New[string]()
	for _, r := range releases {
		releasesSet.Insert(r.Key())
	}
	allHelmResources, err := helm.ListHelmResources(ctx, kc, opts.Namespace)
	if err != nil {
		return nil, err
	}

	leaked := []*common.KubeResource{}
	for _, r := range allHelmResources {
		release := helm.ReleaseFromKubeResource(r)
		if release.Name == "" {
			continue
		}
		if !releasesSet.Has(release.Key()) {
			leaked = append(leaked, r)
		}
	}
	return leaked, nil
}
