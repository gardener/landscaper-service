// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	lsschema "github.com/gardener/landscaper/apis/schema"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// The LandscaperDeployment is created to define a deployment of the Landscaper Service.
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
	// Purpose contains the purpose of this LandscaperDeployment.
	Purpose string `json:"purpose"`
	// Region selects the region this LandscaperDeployment should be installed on.
	// +optional
	Region string `json:"region,omitempty"`
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

var LandscaperDeploymentDefinition = lsschema.CustomResourceDefinition{
	Names: lsschema.CustomResourceDefinitionNames{
		Plural:   "landscaperdeployments",
		Singular: "landscaperdeployment",
		ShortNames: []string{
			"lsdepl",
		},
		Kind: "LandscaperDeployment",
	},
	Scope:             lsschema.NamespaceScoped,
	Storage:           true,
	Served:            true,
	SubresourceStatus: true,
	AdditionalPrinterColumns: []lsschema.CustomResourceColumnDefinition{
		{
			Name:     "Instance",
			Type:     "string",
			JSONPath: ".status.instanceRef.name",
		},
		{
			Name:     "Age",
			Type:     "date",
			JSONPath: ".metadata.creationTimestamp",
		},
	},
}
