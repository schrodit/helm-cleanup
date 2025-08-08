package helm

import (
	"context"
	"fmt"
	"os"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	common "github.com/schrodit/helm-cleanup/pkg/common"
)

func ReleaseFromKubeResource(r *common.KubeResource) common.Release {
	return common.Release{
		Name:      r.GetAnnotations()[HelmReleaseNameAnnotation],
		Namespace: r.GetAnnotations()[HelmReleaseNamespaceAnnotation],
	}
}

func ListHelmResources(ctx context.Context, kc *common.KubeClient, namespace string) ([]*common.KubeResource, error) {
	groups, err := kc.Default.Discovery().ServerPreferredResources()
	if err != nil {
		return nil, errors.Wrap(err, "unable to list available resources")
	}
	resources := []*common.KubeResource{}

	pw := progressbar.Default(int64(len(groups)), "Fetch K8s resources")
	for _, group := range groups {
		if err := pw.Add(1); err != nil {
			return nil, err
		}
		for _, resource := range group.APIResources {
			// Skip subresources like pod/logs, pod/status
			if strings.Contains(resource.Name, "/") {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    group.GroupVersion,
				Version:  resource.Version,
				Resource: resource.Name,
			}
			if gvr.Group == "v1" {
				gvr.Version = gvr.Group
				gvr.Group = ""
			}

			list, err := kc.Dynamic.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s", AppManagedByLabel, AppManagedByHelm),
			})
			if err != nil {
				if k8serrors.IsNotFound(err) || k8serrors.IsMethodNotSupported(err) {
					continue
				}
				return nil, errors.Wrapf(err, "failed to list resource group %q version %q resource %q", gvr.Group, gvr.Version, gvr.Resource)
			}
			for _, item := range list.Items {
				u := item
				resources = append(resources, &common.KubeResource{
					Unstructured:         &u,
					GroupVersionResource: gvr,
				})
			}
		}
	}
	return resources, nil
}

func PrintK8sResourceTable(resources []*common.KubeResource) {
	t := table.NewWriter()
	t.SetStyle(table.StyleDefault)
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{"API VERSION", "KIND", "NAMESPACE", "NAME", "HELM RELEASE"})
	for _, r := range resources {
		release := ReleaseFromKubeResource(r)
		t.AppendRow(table.Row{r.GetAPIVersion(), r.GetKind(), r.GetNamespace(), r.GetName(), release.Key()})
	}
	t.Render()
}
