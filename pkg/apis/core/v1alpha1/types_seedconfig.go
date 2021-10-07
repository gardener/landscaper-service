// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	lsschema "github.com/gardener/landscaper/apis/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// The SeedConfig is created to define the configuration for a Kubernetes cluster, that can host Landscaper Service deployments.
type SeedConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the SeedConfig
	Spec SeedConfigSpec `json:"spec"`

	// Status contains the status of the SeedConfig.
	// +optional
	Status SeedConfigStatus `json:"status"`
}

// SeedConfigSpec contains the specification for a SeedConfig.
type SeedConfigSpec struct {
	// ProviderType specifies the type of the underlying infrastructure provide.
	ProviderType string `json:"providerType"`

	// Region specifies the region in which the target cluster is located.
	Region string `json:"region"`

	// The Priority of this SeedConfig.
	// SeedConfigs with a higher priority number will be preferred over lower numbers
	// when scheduling new landscaper service installations.
	Priority int64 `json:"priority"`

	// Visible defines whether the SeedConfig is visible for scheduling.
	// If set to true, new Landscaper Service deployments can be scheduled on this seed.
	// If set to false, no new Landscaper Service deployments can be scheduled on this seed.
	Visible bool `json:"visible"`

	// SecretRef references the secret that contains the kubeconfig of the target cluster.
	SecretRef ObjectReference `json:"secretRef"`
}

// SeedConfigStatus contains the status of a SeedConfig.
type SeedConfigStatus struct {
	// ObservedGeneration is the most recent generation observed for this SeedConfig.
	// It corresponds to the SeedConfig generation, which is updated on mutation by the landscaper service controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration"`

	// Capacity specifies the remaining capacity for Landscaper Service deployments.
	// For each Landscaper Service deployment that is installed on this seed, the value will be decremented.
	// For each Landscaper Service deployment that is uninstalled from this seed, the value will be incremented.
	// When this value reaches zero, no new Landscaper Services can be deployed on this seed.
	// +optional
	Capacity int64 `json:"capacity"`
}

var SeedConfigDefinition = lsschema.CustomResourceDefinition{
	Names: lsschema.CustomResourceDefinitionNames{
		Plural:   "seedconfigs",
		Singular: "seedconfig",
		ShortNames: []string{
			"seedcfg",
		},
		Kind: "SeedConfig",
	},
	Scope:             lsschema.NamespaceScoped,
	Storage:           true,
	Served:            true,
	SubresourceStatus: true,
	AdditionalPrinterColumns: []lsschema.CustomResourceColumnDefinition{
		{
			Name:     "ProviderType",
			Type:     "string",
			JSONPath: ".spec.providerType",
		},
		{
			Name:     "Region",
			Type:     "string",
			JSONPath: ".spec.region",
		},
		{
			Name:     "Visible",
			Type:     "boolean",
			JSONPath: ".spec.visible",
		},
		{
			Name:     "Priority",
			Type:     "number",
			JSONPath: ".spec.priority",
		},
		{
			Name:     "Capacity",
			Type:     "number",
			JSONPath: ".status.capacity",
		},
		{
			Name:     "Age",
			Type:     "date",
			JSONPath: ".metadata.creationTimestamp",
		},
	},
}
