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

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"

	"k8c.io/kubermatic/v2/pkg/resources/reconciling"

	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func NetworkAttachmentDefinitionCreator(nad *nadv1.NetworkAttachmentDefinition) reconciling.NamedNadV1NetworkAttachmentDefinitionCreatorGetter {
	return func() (string, reconciling.NadV1NetworkAttachmentDefinitionCreator) {
		return nad.Name, func(n *nadv1.NetworkAttachmentDefinition) (*nadv1.NetworkAttachmentDefinition, error) {
			n.ObjectMeta.Annotations = nad.ObjectMeta.Annotations
			n.Spec = nad.Spec
			return n, nil
		}
	}
}

// reconcileNetworkAttachmentDefinition reconciles the needed NetworkAttachmentDefinition.
func reconcileNetworkAttachmentDefinition(ctx context.Context, nad *nadv1.NetworkAttachmentDefinition, namespace string, client ctrlruntimeclient.Client) error {
	creators := []reconciling.NamedNadV1NetworkAttachmentDefinitionCreatorGetter{
		NetworkAttachmentDefinitionCreator(nad),
	}

	if nad != nil {
		if err := reconciling.ReconcileNadV1NetworkAttachmentDefinitions(ctx, creators, namespace, client); err != nil {
			return fmt.Errorf("failed to reconcile NetworkAttachmnetDefinition: %w", err)
		}
	}
	return nil
}
