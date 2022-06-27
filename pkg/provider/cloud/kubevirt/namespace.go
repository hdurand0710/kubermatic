/*
Copyright 2022 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubevirt

import (
	"context"
	"fmt"

	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	kuberneteshelper "k8c.io/kubermatic/v2/pkg/kubernetes"
	"k8c.io/kubermatic/v2/pkg/provider"
	"k8c.io/kubermatic/v2/pkg/resources/reconciling"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	vlanAnnotationKey = "kubermatic.k8c.io/vlan"
	cidrAnnotationKey = "kubermatic.k8c.io/cidr"
)

func NamespaceCreator(name string, cidr, vlan *string) reconciling.NamedNamespaceCreatorGetter {
	return func() (string, reconciling.NamespaceCreator) {
		return name, func(n *corev1.Namespace) (*corev1.Namespace, error) {
			if cidr != nil || vlan != nil {
				if n.Annotations == nil {
					n.Annotations = map[string]string{}
				}
				if vlan != nil {
					n.Annotations[vlanAnnotationKey] = *vlan
				}
				if cidr != nil {
					n.Annotations[cidrAnnotationKey] = *cidr
				}
			}
			return n, nil
		}
	}
}

// reconcileNamespace reconciles a dedicated namespace in the underlying KubeVirt cluster.
func reconcileNamespace(ctx context.Context, name string, cidr, vlan *string, cluster *kubermaticv1.Cluster, update provider.ClusterUpdater, client ctrlruntimeclient.Client) (*kubermaticv1.Cluster, error) {
	creators := []reconciling.NamedNamespaceCreatorGetter{
		NamespaceCreator(name, cidr, vlan),
	}

	if err := reconciling.ReconcileNamespaces(ctx, creators, "", client); err != nil {
		return cluster, fmt.Errorf("failed to reconcile Namespace: %w", err)
	}

	return update(ctx, cluster.Name, func(updatedCluster *kubermaticv1.Cluster) {
		kuberneteshelper.AddFinalizer(updatedCluster, FinalizerNamespace)
	})
}

// deleteNamespace deletes the dedicated namespace.
func deleteNamespace(ctx context.Context, name string, client ctrlruntimeclient.Client) error {
	ns := &corev1.Namespace{}
	if err := client.Get(ctx, types.NamespacedName{Name: name}, ns); err != nil {
		return ctrlruntimeclient.IgnoreNotFound(err)
	}

	return client.Delete(ctx, ns)
}
