// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	LandscaperDeploymentDataPlaneTypeExternal = "External"
	LandscaperDeploymentDataPlaneTypeInternal = "Internal"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LandscaperDeploymentList contains a list of LandscaperDeployment
type LandscaperDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LandscaperDeployment `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// The LandscaperDeployment is created to define a deployment of the landscaper.
// +kubebuilder:resource:singular="landscaperdeployment",path="landscaperdeployments",shortName="lsdepl",scope="Namespaced"
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="DataPlaneType",type=string,JSONPath=`.status.dataPlaneType`
// +kubebuilder:printcolumn:name="Instance",type=string,JSONPath=`.status.instanceRef.name`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
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

	// HighAvailabilityConfig specifies the HA configuration of the resource cluster (shoot cluster)
	// +optional
	HighAvailabilityConfig *HighAvailabilityConfig `json:"highAvailabilityConfig,omitempty"`

	// DataPlane references an externally created and maintained Kubernetes cluster,
	// used as the data plane where Landscaper resources are stored.
	// When DataPlane is defined, the Landscaper Service controller will no longer
	// create its own Kubernetes cluster.
	// +optional
	DataPlane *DataPlane `json:"dataPlane,omitempty"`
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

	// Phase represents the phase of the corresponding Landscaper Instance Installation phase.
	// +optional
	Phase string `json:"phase,omitempty"`

	// DataPlaneType shows whether this deployment has an internal or external data plane cluster.
	// +optional
	DataPlaneType string `json:"dataPlaneType,omitempty"`
}

func (ld *LandscaperDeployment) IsExternalDataPlane() bool {
	return ld.Spec.DataPlane != nil
}

func (ld *LandscaperDeployment) IsInternalDataPlane() bool {
	return ld.Spec.DataPlane == nil
}
