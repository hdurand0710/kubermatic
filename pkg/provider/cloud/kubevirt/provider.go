/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

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
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"

	nadv1 "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"

	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	kuberneteshelper "k8c.io/kubermatic/v2/pkg/kubernetes"
	"k8c.io/kubermatic/v2/pkg/provider"
	"k8c.io/kubermatic/v2/pkg/resources"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// FinalizerNamespace will ensure the deletion of the dedicated namespace.
	FinalizerNamespace = "kubermatic.k8c.io/cleanup-kubevirt-namespace"
	// bridgeName is the OVS bridge name that should exists in each KubeVirt node on the infra cluster.
	bridgeName = "kubevirt"
	// cniVersion is the SemVer of the CNI specification.
	cniVersion                     = "0.4.0"
	defaultCIDRBlock               = "172.24.0.0/16"
	defaultNumberOfAddressePerVlan = 256
)

type kubevirt struct {
	secretKeySelector provider.SecretKeySelectorValueFunc
	dc                kubermaticv1.DatacenterSpecKubevirt
}

func NewCloudProvider(dc *kubermaticv1.Datacenter, secretKeyGetter provider.SecretKeySelectorValueFunc) (provider.CloudProvider, error) {
	if dc.Spec.Kubevirt == nil {
		return nil, errors.New("datacenter is not an KubeVirt datacenter")
	}
	return &kubevirt{
		dc:                *dc.Spec.Kubevirt,
		secretKeySelector: secretKeyGetter,
	}, nil
}

var _ provider.ReconcilingCloudProvider = &kubevirt{}

func (k *kubevirt) DefaultCloudSpec(ctx context.Context, spec *kubermaticv1.CloudSpec) error {
	if spec.Kubevirt != nil {
		return updateInfraStorageClassesInfo(ctx, spec, k.secretKeySelector)
	}
	return nil
}

func (k *kubevirt) ValidateCloudSpec(ctx context.Context, spec kubermaticv1.CloudSpec) error {
	kubeconfig, err := GetCredentialsForCluster(spec, k.secretKeySelector)
	if err != nil {
		return err
	}

	config, err := base64.StdEncoding.DecodeString(kubeconfig)
	if err != nil {
		// if the decoding failed, the kubeconfig is sent already decoded without the need of decoding it,
		// for example the value has been read from Vault during the ci tests, which is saved as json format.
		config = []byte(kubeconfig)
	}

	_, err = clientcmd.RESTConfigFromKubeConfig(config)
	if err != nil {
		return err
	}

	spec.Kubevirt.Kubeconfig = string(config)

	return nil
}

func (k *kubevirt) InitializeCloudProvider(ctx context.Context, cluster *kubermaticv1.Cluster, update provider.ClusterUpdater) (*kubermaticv1.Cluster, error) {
	return k.reconcileCluster(ctx, cluster, update)
}

func (k *kubevirt) ReconcileCluster(ctx context.Context, cluster *kubermaticv1.Cluster, update provider.ClusterUpdater) (*kubermaticv1.Cluster, error) {
	return k.reconcileCluster(ctx, cluster, update)
}

func (k *kubevirt) reconcileCluster(ctx context.Context, cluster *kubermaticv1.Cluster, update provider.ClusterUpdater) (*kubermaticv1.Cluster, error) {
	// Reconcile CSI access: Role and Rolebinding
	client, restConfig, err := k.GetClientWithRestConfigForCluster(cluster)
	if err != nil {
		return cluster, err
	}

	cluster, err = reconcileNamespace(ctx, cluster.Status.NamespaceName, nil, nil, cluster, update, client)
	if err != nil {
		return cluster, err
	}

	err = reconcileCSIRoleRoleBinding(ctx, cluster.Status.NamespaceName, client, restConfig)
	if err != nil {
		return cluster, err
	}

	err = reconcilePresets(ctx, cluster.Status.NamespaceName, client)
	if err != nil {
		return cluster, err
	}

	err = reconcilePreAllocatedDataVolumes(ctx, cluster, client)
	if err != nil {
		return cluster, err
	}

	// Reconcile NetworkAttachmentDefinition only if option is enabled in Cluster
	// VlanIsolationEnabled is immutable, it's not possible to go from !=nil to nil (or true to false), then there is intentionally no cleanup
	if cluster.Spec.Cloud.Kubevirt.VlanIsolation != nil && cluster.Spec.Cloud.Kubevirt.VlanIsolation.Enabled {
		allNs := &corev1.NamespaceList{}
		selector := fields.OneTermNotEqualSelector("metadata.name", cluster.Status.NamespaceName) // excludes cluster ns itself
		options := &ctrlruntimeclient.ListOptions{FieldSelector: selector}
		if err := client.List(ctx, allNs, options); err != nil {
			return nil, err
		}

		vlan := k.findAvailableVlan(allNs)
		cidr, err := k.findAvailableCidr(allNs)
		if err != nil {
			return cluster, err
		}

		config := fmt.Sprintf(`
		{
			"cniVersion": "%s",
			"type": "ovs", 
			"bridge": "%s", 
			"vlan": %s, 
			"ipam": {
				"type": "whereabouts", 
				"range": "%s"
			}
		}
		`, cniVersion, bridgeName, vlan, cidr)

		nad := &nadv1.NetworkAttachmentDefinition{
			ObjectMeta: metav1.ObjectMeta{
				Name:      cluster.Spec.Cloud.Kubevirt.VlanIsolation.NetworkName,
				Namespace: cluster.Status.NamespaceName,
				Annotations: map[string]string{
					"k8s.v1.cni.cncf.io/resourceName": fmt.Sprintf("ovs-cni.network.kubevirt.io/%s", bridgeName),
				},
			},
			Spec: nadv1.NetworkAttachmentDefinitionSpec{
				Config: config,
			},
		}
		err = reconcileNetworkAttachmentDefinition(ctx, nad, cluster.Status.NamespaceName, client)
		if err != nil {
			return cluster, err
		}

		cluster, err = reconcileNamespace(ctx, cluster.Status.NamespaceName, &cidr, &vlan, cluster, update, client)
	}

	return cluster, err
}

// findAvailableVlan iterates over all namespaces in the KubeVirt infra cluster to get the annotation:
// // - vlanAnnotationKey
// Finds the first available and returns it.
func (k *kubevirt) findAvailableVlan(allNs *corev1.NamespaceList) string {
	existingVlans := getExistingAnnotationValues(allNs, vlanAnnotationKey)

	// Find first available
	i := 1
	for {
		v := strconv.Itoa(i)
		if _, ok := existingVlans[v]; !ok {
			return v
		}
		i++
	}
}

// findAvailableCidr iterates over all namespaces in the KubeVirt infra cluster to get the annotation:
// - cidrAnnotationKey
// Finds the first available and returns it.
func (k *kubevirt) findAvailableCidr(allNs *corev1.NamespaceList) (string, error) {
	existingCidrs := getExistingAnnotationValues(allNs, cidrAnnotationKey)

	// cidr
	// All possible ranges from DC settings
	ranges := k.possibleRanges()
	numberOfAddressesPerVlan := numberOfAddressesPerVlan(k.dc)

	// split them into all possible vlan sub ranges
	vlanCIDRs, err := ranges.SplitByNumberOfAddresses(numberOfAddressesPerVlan)
	if err != nil {
		return "", err
	}

	// Find the first vlan CIDR available
	var cidr string
	for _, c := range vlanCIDRs.CIDRBlocks {
		if _, ok := existingCidrs[c]; !ok {
			return c, nil
		}
	}

	return cidr, nil
}

func (k *kubevirt) possibleRanges() *kubermaticv1.NetworkRanges {
	ranges := kubermaticv1.NetworkRanges{}
	ranges.CIDRBlocks = make([]string, 0)
	for _, r := range k.dc.VlanIsolation.VlanRanges {
		ranges.CIDRBlocks = append(ranges.CIDRBlocks, string(r))
	}

	// Default values, should not be used, as their should be a default value in the Seed CR
	if len(ranges.CIDRBlocks) == 0 {
		ranges.CIDRBlocks = append(ranges.CIDRBlocks, defaultCIDRBlock)
	}

	return &ranges
}

func getExistingAnnotationValues(allNs *corev1.NamespaceList, annotationKey string) map[string]bool {
	values := map[string]bool{}

	for _, ns := range allNs.Items {
		if vlan, ok := ns.Annotations[annotationKey]; ok {
			values[vlan] = true
		}
	}
	return values
}

func numberOfAddressesPerVlan(dc kubermaticv1.DatacenterSpecKubevirt) (nb uint32) {
	nb = dc.VlanIsolation.NumberOfAddressesPerVlan

	// Default values, should not be used, as their should be a default value in the Seed CR
	if nb == 0 {
		nb = defaultNumberOfAddressePerVlan
	}
	return
}

func (k *kubevirt) CleanUpCloudProvider(ctx context.Context, cluster *kubermaticv1.Cluster, update provider.ClusterUpdater) (*kubermaticv1.Cluster, error) {
	client, _, err := k.GetClientWithRestConfigForCluster(cluster)
	if err != nil {
		return cluster, err
	}

	if kuberneteshelper.HasFinalizer(cluster, FinalizerNamespace) {
		if err := deleteNamespace(ctx, cluster.Status.NamespaceName, client); err != nil && !apierrors.IsNotFound(err) {
			return cluster, fmt.Errorf("failed to delete namespace %s: %w", cluster.Status.NamespaceName, err)
		}
		return update(ctx, cluster.Name, func(updatedCluster *kubermaticv1.Cluster) {
			kuberneteshelper.RemoveFinalizer(updatedCluster, FinalizerNamespace)
		})
	}

	return cluster, nil
}

func (k *kubevirt) ValidateCloudSpecUpdate(ctx context.Context, oldSpec kubermaticv1.CloudSpec, newSpec kubermaticv1.CloudSpec) error {
	// Immutable fields
	if oldSpec.Kubevirt.VlanIsolation != nil && oldSpec.Kubevirt.VlanIsolation != newSpec.Kubevirt.VlanIsolation {
		return fmt.Errorf("updating KubeVirt Vlan Isolation (was %v, updated to %v)", oldSpec.Kubevirt.VlanIsolation, newSpec.Kubevirt.VlanIsolation)
	}

	return nil
}

// GetClientWithRestConfigForCluster returns the kubernetes client and the rest config for the KubeVirt underlying cluster.
func (k *kubevirt) GetClientWithRestConfigForCluster(cluster *kubermaticv1.Cluster) (ctrlruntimeclient.Client, *restclient.Config, error) {
	if cluster.Spec.Cloud.Kubevirt == nil {
		return nil, nil, errors.New("No KubeVirt provider spec")
	}
	kubeconfig, err := GetCredentialsForCluster(cluster.Spec.Cloud, k.secretKeySelector)
	if err != nil {
		return nil, nil, err
	}

	client, restConfig, err := NewClientWithRestConfig(kubeconfig)
	if err != nil {
		return nil, nil, err
	}

	if err := kubevirtv1.AddToScheme(client.Scheme()); err != nil {
		return nil, nil, err
	}
	if err = cdiv1beta1.AddToScheme(client.Scheme()); err != nil {
		return nil, nil, err
	}
	if err = nadv1.AddToScheme(client.Scheme()); err != nil {
		return nil, nil, err
	}

	return client, restConfig, nil
}

// GetCredentialsForCluster returns the credentials for the passed in cloud spec or an error.
func GetCredentialsForCluster(cloud kubermaticv1.CloudSpec, secretKeySelector provider.SecretKeySelectorValueFunc) (kubeconfig string, err error) {
	kubeconfig = cloud.Kubevirt.Kubeconfig

	if kubeconfig == "" {
		if cloud.Kubevirt.CredentialsReference == nil {
			return "", errors.New("no credentials provided")
		}
		kubeconfig, err = secretKeySelector(cloud.Kubevirt.CredentialsReference, resources.KubevirtKubeConfig)
		if err != nil {
			return "", err
		}
	}

	return kubeconfig, nil
}
