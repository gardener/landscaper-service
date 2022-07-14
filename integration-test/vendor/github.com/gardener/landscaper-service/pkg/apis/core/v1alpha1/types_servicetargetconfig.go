// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	lsschema "github.com/gardener/landscaper/apis/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

var ServiceTargetConfigDefinition = lsschema.CustomResourceDefinition{
	Names: lsschema.CustomResourceDefinitionNames{
		Plural:   "servicetargetconfigs",
		Singular: "servicetargetconfig",
		ShortNames: []string{
			"servcfg",
		},
		Kind: "ServiceTargetConfig",
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
			JSONPath: ".metadata.labels.config\\.landscaper-service\\.gardener\\.cloud/region",
		},
		{
			Name:     "Visible",
			Type:     "string",
			JSONPath: ".metadata.labels.config\\.landscaper-service\\.gardener\\.cloud/visible",
		},
		{
			Name:     "Priority",
			Type:     "number",
			JSONPath: ".spec.priority",
		},
		{
			Name:     "Age",
			Type:     "date",
			JSONPath: ".metadata.creationTimestamp",
		},
	},
}
