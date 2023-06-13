// SPDX-FileCopyrightText: 2023 "SAP SE or an SAP affiliate company and Gardener contributors"
//
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	lsschema "github.com/gardener/landscaper/apis/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TenantRegistrationList contains a list of TenantRegistration
type TenantRegistrationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TenantRegistration `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type TenantRegistration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec contains the specification for the TenantRegistration.
	Spec TenantRegistrationSpec `json:"spec"`

	// Status contains the status for the TenantRegistration.
	// +optional
	Status TenantRegistrationStatus `json:"status"`
}

// TenantRegistrationStatus contains the status for the TenantRegistration.
type TenantRegistrationStatus struct {
	SyncedGeneration   int64  `json:"syncedGeneration"`
	ObservedGeneration int64  `json:"observedGeneration"`
	Namespace          string `json:"namespace"`
}

// TenantRegistrationSpec contains the specification for the TenantRegistration.
type TenantRegistrationSpec struct {
	Author string `json:"author"`
}

var TenantRegistrationDefinition = lsschema.CustomResourceDefinition{
	Names: lsschema.CustomResourceDefinitionNames{
		Plural:   "tenantregistrations",
		Singular: "tenantregistration",
		ShortNames: []string{
			"ten",
			"tenreg",
		},
		Kind: "TenantRegistration",
	},
	Scope:             lsschema.ClusterScoped,
	Storage:           true,
	Served:            true,
	SubresourceStatus: true,
}
