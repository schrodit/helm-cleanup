package cleanup

import (
	"context"
	"log"

	"github.com/schrodit/helm-cleanup/pkg/common"
	"github.com/schrodit/helm-cleanup/pkg/helm"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/sets"
)

func ListLeakedResources(ctx context.Context, opts common.Options) ([]*unstructured.Unstructured, error) {
	releases, err := helm.ListReleases(opts)
	if err != nil {
		return nil, err
	}
	if opts.DryRun || opts.Debug {
		log.Println("RELEASES")
	}
	releasesSet := sets.New[string]()
	for _, r := range releases {
		if opts.DryRun || opts.Debug {
			log.Printf("%s/%s\n", r.Namespace, r.Name)
		}
		releasesSet.Insert(r.Key())
	}
	allHelmResources, err := helm.ListHelmResources(ctx, opts.KubeConfig, opts.Namespace)
	if err != nil {
		return nil, err
	}

	leaked := []*unstructured.Unstructured{}
	for _, u := range allHelmResources {
		release := helm.ReleaseFromUnstructured(&u)
		if (release.Name == "") {
			continue
		}
		if !releasesSet.Has(release.Key()) {
			leaked = append(leaked, u.DeepCopy())
		}
	}
	return leaked, nil
}
