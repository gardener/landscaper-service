// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package core

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LandscaperDeploymentList contains a list of LandscaperDeployment
type LandscaperDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LandscaperDeployment `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// The LandscaperDeployment is created to define a deployment of the landscaper.
type LandscaperDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the LandscaperDeployment
	Spec LandscaperDeploymentSpec `json:"spec"`

	// Status contains the status of the LandscaperDeployment.
	// +optional
	Status LandscaperDeploymentStatus `json:"status"`
}

// LandscaperDeploymentSpec contains the specification for a LandscaperDeployment.
type LandscaperDeploymentSpec struct {
	// TenantId is the unique identifier of the owning tenant.
	TenantId string `json:"tenantId"`

	// Purpose contains the purpose of this LandscaperDeployment.
	Purpose string `json:"purpose"`

	// LandscaperConfiguration contains the configuration for the landscaper service deployment
	LandscaperConfiguration LandscaperConfiguration `json:"landscaperConfiguration"`

	// OIDCConfig describes the OIDC config of the customer resource cluster (shoot cluster)
	// +optional
	OIDCConfig *OIDCConfig `json:"oidcConfig,omitempty"`
}

// LandscaperDeploymentStatus contains the status of a LandscaperDeployment.
type LandscaperDeploymentStatus struct {
	// ObservedGeneration is the most recent generation observed for this LandscaperDeployment.
	// It corresponds to the LandscaperDeployment generation, which is updated on mutation by the landscaper service controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration"`

	// LastError describes the last error that occurred.
	// +optional
	LastError *Error `json:"lastError,omitempty"`

	// InstanceRef references the instance that is created for this LandscaperDeployment.
	// +optional
	InstanceRef *ObjectReference `json:"instanceRef"`
}
