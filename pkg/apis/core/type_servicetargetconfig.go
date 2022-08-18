// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ServiceTargetConfigVisibleLabelName label defines whether the ServiceTargetConfig is visible for scheduling.
	// If set to "true", any Landscaper Service deployment can be scheduled on this seed.
	// If not set or set to "false", no new Landscaper Service deployments can be scheduled on this seed.
	ServiceTargetConfigVisibleLabelName = "config.landscaper-service.gardener.cloud/visible"
	// ServiceTargetConfigRegionLabelName label specifies the region in which the target cluster is located.
	ServiceTargetConfigRegionLabelName = "config.landscaper-service.gardener.cloud/region"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceTargetConfigList contains a list of ServiceTargetConfig
type ServiceTargetConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceTargetConfig `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// The ServiceTargetConfig is created to define the configuration for a Kubernetes cluster, that can host Landscaper Service deployments.
type ServiceTargetConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the ServiceTargetConfig
	Spec ServiceTargetConfigSpec `json:"spec"`

	// Status contains the status of the ServiceTargetConfig.
	// +optional
	Status ServiceTargetConfigStatus `json:"status"`
}

// ServiceTargetConfigSpec contains the specification for a ServiceTargetConfig.
type ServiceTargetConfigSpec struct {
	// ProviderType specifies the type of the underlying infrastructure provide.
	ProviderType string `json:"providerType"`

	// The Priority of this ServiceTargetConfig.
	// SeedConfigs with a higher priority number will be preferred over lower numbers
	// when scheduling new landscaper service installations.
	Priority int64 `json:"priority"`

	// SecretRef references the secret that contains the kubeconfig of the target cluster.
	SecretRef SecretReference `json:"secretRef"`
}

// ServiceTargetConfigStatus contains the status of a ServiceTargetConfig.
type ServiceTargetConfigStatus struct {
	// ObservedGeneration is the most recent generation observed for this ServiceTargetConfig.
	// It corresponds to the ServiceTargetConfig generation, which is updated on mutation by the landscaper service controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration"`

	// InstanceRefs is the list of references to instances that use this ServiceTargetConfig.
	// +optional
	InstanceRefs []ObjectReference `json:"instanceRefs,omitempty"`
}
