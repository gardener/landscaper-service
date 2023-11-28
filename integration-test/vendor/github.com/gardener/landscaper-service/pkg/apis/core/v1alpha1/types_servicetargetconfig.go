// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceTargetConfigList contains a list of ServiceTargetConfig
type ServiceTargetConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServiceTargetConfig `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// The ServiceTargetConfig is created to define the configuration for a Kubernetes cluster, that can host Landscaper Service deployments.
// +kubebuilder:resource:singular="servicetargetconfig",path="servicetargetconfigs",shortName="servcfg",scope="Namespaced"
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Visible",type=string,JSONPath=`.metadata.labels.config\.landscaper-service\.gardener\.cloud/visible`
// +kubebuilder:printcolumn:name="Priority",type=number,JSONPath=`.spec.priority`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
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

	// The Priority of this ServiceTargetConfig.
	// SeedConfigs with a higher priority number will be preferred over lower numbers
	// when scheduling new landscaper service installations.
	Priority int64 `json:"priority"`

	// SecretRef references the secret that contains the kubeconfig of the target cluster.
	SecretRef SecretReference `json:"secretRef"`

	// IngressDomain is the ingress domain of the corresponding target cluster.
	IngressDomain string `json:"ingressDomain"`
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
