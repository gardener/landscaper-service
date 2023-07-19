// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0
package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstanceRegistrationList contains a list of InstanceRegistration
type InstanceRegistrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InstanceRegistration `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// The InstanceRegistration is created by the tenant and will be translated to a LandscaperDeployment.
type InstanceRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the InstanceRegistration.
	Spec InstanceRegistrationSpec `json:"spec"`

	// Status contains the status for the InstanceRegistration.
	// +optional
	Status InstanceRegistrationStatus `json:"status"`
}

type InstanceRegistrationSpec struct {
	LandscaperDeploymentSpec
}

type InstanceRegistrationStatus struct {
	// ObservedGeneration is the most recent generation observed for this InstanceRegistration.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration"`

	// LandscaperDeploymentInfo contains the namespace/name for the corresponding LandscaperDeployment CR
	LandscaperDeploymentInfo *LandscaperDeploymentInfo `json:"landscaperDeployment,omitempty"`

	// LastError describes the last error that occurred.
	// +optional
	LastError *Error `json:"lastError,omitempty"`

	// UserKubeconfig contains the user kubeconfig which can be used for accessing the landscaper cluster.
	// +optional
	UserKubeconfig string `json:"userKubeconfig,omitempty"`
}

type LandscaperDeploymentInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
