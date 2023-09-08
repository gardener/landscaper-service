// SPDX-FileCopyrightText: 2021 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	lsschema "github.com/gardener/landscaper/apis/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceRegistrationList contains a list of NamespaceRegistration
type NamespaceRegistrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NamespaceRegistration `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NamespaceRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the NamespaceRegistration.
	Spec NamespaceRegistrationSpec `json:"spec"`

	// Status contains the status for the NamespaceRegistration.
	// +optional
	Status NamespaceRegistrationStatus `json:"status"`
}

type NamespaceRegistrationStatus struct {
	Phase          string      `json:"phase"`
	Description    string      `json:"description"`
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
	// +optional
	LastError *Error `json:"lastError,omitempty"`
}

type NamespaceRegistrationSpec struct {
}

var NamespaceRegistrationDefinition = lsschema.CustomResourceDefinition{
	Names: lsschema.CustomResourceDefinitionNames{
		Plural:   "namespaceregistrations",
		Singular: "namespaceregistration",
		ShortNames: []string{
			"nsreg",
		},
		Kind: "NamespaceRegistration",
	},
	Scope:             lsschema.NamespaceScoped,
	Storage:           true,
	Served:            true,
	SubresourceStatus: true,
	AdditionalPrinterColumns: []lsschema.CustomResourceColumnDefinition{
		{
			Name:     "Phase",
			Type:     "string",
			JSONPath: ".status.phase",
		},
	},
}
