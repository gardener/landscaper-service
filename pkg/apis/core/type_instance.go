// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// The Instance is created for each LandscaperDeployment.
// The landscaper service controller selects a suitable/available SeedConfig and creates
// an Installation.
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the Instance.
	Spec InstanceSpec `json:"spec"`

	// Status contains the status for the Instance.
	// +optional
	Status InstanceStatus `json:"status"`
}

// InstanceSpec contains the specification for an Instance.
type InstanceSpec struct {
	// SeedConfigReg specifies the target cluster for which the installation is created.
	SeedConfigRef ObjectReference `json:"seedConfigRef"`
}

// InstanceStatus contains the status for an Instance.
type InstanceStatus struct {
	// ObservedGeneration is the most recent generation observed for this Instance.
	// It corresponds to the Instance generation, which is updated on mutation by the landscaper service controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration"`

	// LastError describes the last error that occurred.
	// +optional
	LastError *Error `json:"lastError,omitempty"`

	// TargetRef references the Target for this Instance.
	// +optional
	TargetRef *ObjectReference `json:"targetRef"`

	// InstallationRef references the Installation for this Instance.
	// +optional
	InstallationRef *ObjectReference `json:"installationRef"`

	// ClusterEndpointRef references a data object,
	// containing the URL at which the landscaper cluster is accessible.
	// +optional
	ClusterEndpointRef *ObjectReference `json:"clusterEndpoint"`

	// ClusterEndpointRef references a data object,
	// containing the Kubeconfig which can be used for accessing the landscaper cluster.
	// +optional
	ClusterKubeconfigRef *ObjectReference `json:"clusterKubeconfigRef"`
}
