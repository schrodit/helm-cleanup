package helm

import (
	"context"
	"fmt"
	"log"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/pkg/errors"
	common "github.com/schrodit/helm-cleanup/pkg/common"
	"github.com/schrodit/helm-cleanup/pkg/k8s"
)

func ReleaseFromUnstructured(u *unstructured.Unstructured) common.Release {
	return common.Release{
		Name:      u.GetAnnotations()[HelmReleaseNameAnnotation],
		Namespace: u.GetAnnotations()[HelmReleaseNamespaceAnnotation],
	}
}

func ListHelmResources(ctx context.Context, kc common.KubeConfig, namespace string) ([]unstructured.Unstructured, error) {
	client, dynClient, err := k8s.GetClientSetWithKubeConfig(kc.File, kc.Context)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create kubernetes client")
	}

	groups, err := client.Discovery().ServerPreferredResources()
	if err != nil {
		return nil, errors.Wrap(err, "unable to list available resources")
	}

	log.Printf("Found %d groups\n", len(groups))

	resources := []unstructured.Unstructured{}

	for _, group := range groups {
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

			list, err := dynClient.Resource(gvr).Namespace(namespace).List(ctx, metav1.ListOptions{
				LabelSelector: fmt.Sprintf("%s=%s", AppManagedByLabel, AppManagedByHelm),
			})
			if err != nil {
				if k8serrors.IsNotFound(err) || k8serrors.IsMethodNotSupported(err) {
					continue
				}
				return nil, errors.Wrapf(err, "failed to list resource group %q version %q resource %q", gvr.Group, gvr.Version, gvr.Resource)
			}

			resources = append(resources, list.Items...)
		}
	}
	return resources, nil
}
